package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/txn2/kubefwd/pkg/fwdport"
	"github.com/txn2/kubefwd/pkg/utils"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube/servicecontext"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	frpmodels "github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/kubefwd"
	"github.com/we-dcode/kube-tunnel/pkg/models"
	"github.com/we-dcode/kube-tunnel/pkg/utils/hostsutils"
	"github.com/we-dcode/kube-tunnel/pkg/utils/maputils"
	"os"
)

type KubeTunnel struct {
	kubeClient *kube.Kube
	helmClient *helm.Helm
}

type KubeTunnelConf struct {
	ServiceName             string
	KubeTunnelPortMap       map[string]string
	LocalIP                 string
	DnsForwardAllNamespaces bool
	Proxies                 map[string]int
}

func MustNewKubeTunnel(kubeConfig string, namespace string, privileged bool) *KubeTunnel {

	if len(kubeConfig) > 0 {
		_ = os.Setenv("KUBECONFIG", kubeConfig) // Workaround for internal client created by kubefwd
	}

	kubeClient := kube.MustNew(kubeConfig, namespace)

	err := kubeClient.ConnectivityCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	err = kubeClient.RBACCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	if privileged {
		err = CheckRootPermissions()
		if err != nil {
			log.Panicf(err.Error())
		}
	}

	helmClient := helm.MustNew(kubeClient)

	return &KubeTunnel{
		kubeClient: kubeClient,
		helmClient: helmClient,
	}
}

func (ct *KubeTunnel) Install(operatorVersion string) {

	err := ct.helmClient.InstallKubeTunnelOperator(operatorVersion)
	if err != nil {
		log.Panic(err.Error())
	}
}

func (ct *KubeTunnel) CreateTunnel(tunnelConf KubeTunnelConf) {

	svcCtx, err := ct.kubeClient.GetServiceContext(tunnelConf.ServiceName)
	if err != nil {
		log.Panic(err.Error())
	}

	labels, err := ct.kubeClient.GetPodLabelsByLabelSelector(ct.kubeClient.Namespace, svcCtx.LabelSelector)
	if err != nil {
		log.Panic(err.Error())
	}

	kubeTunnelResourceSpec := servicecontext.ToKubeTunnelResourceSpec(svcCtx, labels, tunnelConf.Proxies)

	if err = ct.kubeClient.CreateKubeTunnelResource(kubeTunnelResourceSpec); err != nil {

		log.Panicf("fail creating kubetunnel CRD, internal error: %s", err.Error())
	}

	kubefwdSyncChannel := make(chan error)
	var hostFile *fwdport.HostFileWithLock

	go func(ktrs *models.KubeTunnelResourceSpec) {
		hostFile = kubefwd.Execute(ct.kubeClient, ktrs, kubefwdSyncChannel, tunnelConf.DnsForwardAllNamespaces)
	}(&kubeTunnelResourceSpec)

	err = <-kubefwdSyncChannel

	if err != nil {
		log.Panicf("fail executing kubefwd: %s", err.Error())
	}

	log.Info("executing frpc")

	common := frpmodels.Common{
		ServerAddress: fmt.Sprintf("%s-%s", constants.KubetunnelSlug, kubeTunnelResourceSpec.ServiceName),
		ServerPort:    constants.FRPServerPort,
	}

	servicePortsPairs := servicecontext.ToFRPClientPairs(tunnelConf.LocalIP, tunnelConf.KubeTunnelPortMap, svcCtx)

	frpcManager := frpc.NewManager(common, servicePortsPairs, hostFile)

	if ok, localIP, _ := hostFile.Hosts.HostAddressLookup(common.ServerAddress); ok {

		proxyDNSim := maputils.Keys(tunnelConf.Proxies)
		proxyDNSim = append(proxyDNSim, fmt.Sprintf("%s.%s", proxyDNSim[0], constants.KubetunnelSlug))
		hostFile.Hosts.AddHosts(localIP, proxyDNSim)
		hostFile.Hosts.Save()
	}

	frpcManager.RunFRPc()
	hostsutils.HostsCleanup(hostFile.Hosts)
	hostFile.Hosts.Save()
}

func CheckRootPermissions() error {

	hasRoot, err := utils.CheckRoot()

	if !hasRoot {
		log.Errorf(`
This program requires superuser privileges to run. These
privileges are required to add IP address aliases to your
loopback interface. Superuser privileges are also needed
to listen on low port numbers for these IP addresses.

Try:
 - sudo -E kubetunnel (Unix)
 - Running a shell with administrator rights (Windows)

`)
		if err != nil {
			return fmt.Errorf("root check failure: %s", err.Error())
		}
	}

	return nil
}
