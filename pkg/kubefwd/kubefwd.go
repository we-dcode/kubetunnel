package kubefwd

import (
	"bytes"
	"context"
	"fmt"
	"github.com/txn2/kubefwd/cmd/kubefwd/services"
	"github.com/txn2/kubefwd/pkg/fwdcfg"
	"github.com/txn2/kubefwd/pkg/fwdhost"
	"github.com/txn2/kubefwd/pkg/fwdport"
	"github.com/txn2/kubefwd/pkg/fwdsvcregistry"
	"github.com/txn2/txeh"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/models"
	"github.com/we-dcode/kube-tunnel/pkg/notify/killsignal"
	"github.com/we-dcode/kube-tunnel/pkg/utils/hostsutils"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	"io/ioutil"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/tools/cache"
	"os"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var globalUsage = ``
var Version = "0.0.0"

func init() {
	// quiet version
	args := os.Args[1:]
	if len(args) == 2 && args[0] == "version" && args[1] == "quiet" {
		fmt.Println(Version)
		os.Exit(0)
	}

	log.SetOutput(&LogOutputSplitter{})
	if len(args) > 0 && args[0] == "completion" {
		log.SetOutput(ioutil.Discard)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "kubefwd",
		Short: "Expose Kubernetes services for local development.",
		Example: " kubefwd services --help\n" +
			"  kubefwd svc -n the-project\n" +
			"  kubefwd svc -n the-project -l env=dev,component=api\n" +
			"  kubefwd svc -n the-project -f metadata.name=service-name\n" +
			"  kubefwd svc -n default -l \"app in (ws, api)\"\n" +
			"  kubefwd svc -n default -n the-project\n" +
			"  kubefwd svc -n the-project -m 80:8080 -m 443:1443\n" +
			"  kubefwd svc -n the-project -z path/to/conf.yml\n" +
			"  kubefwd svc -n the-project -r svc.ns:127.3.3.1\n" +
			"  kubefwd svc --all-namespaces",

		Long: globalUsage,
	}

	cmd.AddCommand(services.Cmd)

	return cmd
}

type LogOutputSplitter struct{}

func (splitter *LogOutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=error")) || bytes.Contains(p, []byte("level=warn")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}

// Execute - This code was copied from kubefwd and modified a bit to support kubetunnel requirements
func Execute(kubeClient *kube.Kube, frpsValues *models.KubeTunnelResourceSpec, channel chan error) *fwdport.HostFileWithLock {

	log.Println("Press [Ctrl-C] to stop forwarding.")
	log.Println("'cat /etc/hosts' to see all host entries.")

	hostFile, err := txeh.NewHostsDefault()
	if err != nil {
		log.Fatalf("HostFile error: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("Loaded hosts file %s\n", hostFile.ReadFilePath)

	msg, err := fwdhost.BackupHostFile(hostFile)
	if err != nil {
		log.Fatalf("Error backing up hostfile: %s\n", err.Error())
		os.Exit(1)
	}

	hostsutils.HostsCleanup(hostFile)

	log.Printf("HostFile management: %s", msg)

	// if no context override
	fwdsvcregistry.Init(killsignal.CancellationChannel.C)

	nsWatchesDone := &sync.WaitGroup{} // We'll wait on this to exit the program. Done() indicates that all namespace watches have shutdown cleanly.

	nsWatchesDone.Add(1)

	configGetter := fwdcfg.NewConfigGetter()

	restClient, _ := configGetter.GetRESTClient()

	nameSpaceOpts := services.NamespaceOpts{
		ClientSet: *kubeClient.InnerKubeClient,
		Namespace: kubeClient.Namespace,

		// For parallelization of ip handout,
		// each cluster and namespace has its own ip range
		NamespaceIPLock:   &sync.Mutex{},
		ListOptions:       metav1.ListOptions{},
		HostFile:          &fwdport.HostFileWithLock{Hosts: hostFile},
		ClientConfig:      *kubeClient.Config,
		Domain:            constants.KubetunnelSlug,
		RESTClient:        *restClient,
		ClusterN:          0,
		NamespaceN:        0,
		ManualStopChannel: killsignal.CancellationChannel.C,
	}

	go func(npo services.NamespaceOpts) {
		watchServiceEvents(&nameSpaceOpts, killsignal.CancellationChannel.C)
		nsWatchesDone.Done()
	}(nameSpaceOpts)

	go func() {
		nsWatchesDone.Wait()
		log.Debugf("namespace watchers is done")

		// Shutdown all active services
		<-fwdsvcregistry.Done()

		log.Infof("Clean exit")
	}()

	go func() {
		WaitUntilKubeTunnelIsUp(frpsValues, fwdsvcregistry.Done())
		channel <- nil
	}()

	return nameSpaceOpts.HostFile
}

func WaitUntilKubeTunnelIsUp(frpsValues *models.KubeTunnelResourceSpec, done <-chan struct{}) {

	host := frpsValues.KubeTunnelServiceName()

	// TODO: Do I need to wait for interrupt here or it's already handled?
	for tcputil.IsAvailable(host, constants.FRPServerPort) == false {

		// If Done already request (Interrupt event) then break the loop
		if IsChannelClosed(done) {
			break
		}

		time.Sleep(200 * time.Millisecond)
	}
}

func IsChannelClosed(ch <-chan struct{}) bool {
	select {
	case <-ch:
		return true
	default:
		return false
	}
}

func watchServiceEvents(opts *services.NamespaceOpts, stopListenCh <-chan struct{}) {
	// Apply filtering
	optionsModifier := func(options *metav1.ListOptions) {
		options.FieldSelector = opts.ListOptions.FieldSelector
		options.LabelSelector = opts.ListOptions.LabelSelector
	}

	// Construct the informer object which will query the api server,
	// and send events to our handler functions
	// https://engineering.bitnami.com/articles/kubewatch-an-example-of-kubernetes-custom-controller.html
	_, controller := cache.NewInformer(
		&cache.ListWatch{
			ListFunc: func(options metav1.ListOptions) (runtime.Object, error) {
				optionsModifier(&options)
				return opts.ClientSet.CoreV1().Services(opts.Namespace).List(context.TODO(), options)
			},
			WatchFunc: func(options metav1.ListOptions) (watch.Interface, error) {
				options.Watch = true
				optionsModifier(&options)
				return opts.ClientSet.CoreV1().Services(opts.Namespace).Watch(context.TODO(), options)
			},
		},
		&v1.Service{},
		0,
		cache.ResourceEventHandlerFuncs{
			AddFunc:    opts.AddServiceHandler,
			DeleteFunc: opts.DeleteServiceHandler,
			UpdateFunc: func(oldObj, newObj interface{}) {
				UpdateServiceHandler(opts, oldObj, newObj)
			},
		},
	)

	// Start the informer, blocking call until we receive a stop signal
	controller.Run(stopListenCh)
	log.Infof("Stopped watching Service events in namespace %s in %s context", opts.Namespace, opts.Context)
}

func UpdateServiceHandler(opts *services.NamespaceOpts, oldObj interface{}, newObj interface{}) {
	key, err := cache.MetaNamespaceKeyFunc(newObj)
	if err == nil {
		log.Printf("update service %s. replacing dns port-forwarding", key)
	}

	opts.DeleteServiceHandler(oldObj)
	opts.AddServiceHandler(newObj)
}
