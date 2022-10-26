package kube

import (
	"context"
	"encoding/json"
	"fmt"
	portforward "github.com/maordavidov/go-k8s-portforward"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube/servicecontext"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/models"
	v12 "k8s.io/api/authorization/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"strconv"
	"strings"
)

type Kube struct {
	InnerKubeClient *kubernetes.Clientset
	Config          *rest.Config
	Namespace       string
}

func MustNew(kubeConf string, namespace string) *Kube {

	kube, err := New(kubeConf, namespace)

	if err != nil {
		log.Panic(err.Error())
	}

	return kube
}

func New(kubeConf string, namespace string) (kube *Kube, err error) {

	kubeClient, config, err := createInClusterKubeClient()

	if err != nil {

		var kubeConfig clientcmd.ClientConfig
		kubeClient, kubeConfig, _ = mustCreateOutOfClusterKubeClient(kubeConf)

		config, _ = kubeConfig.ClientConfig()

		if len(namespace) == 0 {
			namespace, _, _ = kubeConfig.Namespace()
		}
	}

	kube = &Kube{
		kubeClient,
		config,
		namespace,
	}

	return kube, nil
}

func (k *Kube) GetPodLabelsByLabelSelector(namespace string, podLabelSelectors map[string]string) (map[string]string, error) {

	labelSelector, err := labels.ValidatedSelectorFromSet(podLabelSelectors)
	if err != nil {
		return nil, err
	}

	pods, err := k.InnerKubeClient.CoreV1().Pods(namespace).List(context.Background(), v1.ListOptions{
		LabelSelector: labelSelector.String(),
		Limit:         1,
	})

	if err != nil {
		return nil, err
	}

	if len(pods.Items) == 0 {

		return nil, fmt.Errorf("could not find pods. host: '%s' namespace: '%s', selector: '%s'", k.Config.Host, namespace, labelSelector.String())
	}

	return pods.Items[0].Labels, nil
}

func (k *Kube) GetServiceContextWithNamespace(name string, namespace string) (*servicecontext.ServiceContext, error) {
	if len(namespace) == 0 {
		namespace = k.Namespace
	}

	svc, err := k.InnerKubeClient.CoreV1().Services(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {

		err = fmt.Errorf("namespace: '%s' svc: '%s' not found at host: '%s'", k.Namespace, name, k.Config.Host)
		return nil, err
	}

	ctx := servicecontext.ServiceContext{
		ServiceName:   svc.Name,
		LabelSelector: svc.Spec.Selector,
		Ports:         svc.Spec.Ports,
	}

	return &ctx, nil
}
func (k *Kube) GetServiceContext(name string) (*servicecontext.ServiceContext, error) {

	return k.GetServiceContextWithNamespace(name, k.Namespace)
}

func (k *Kube) ListServiceNamesWithoutKubeTunnel() (serviceNames []string, err error) {

	services, err := k.InnerKubeClient.CoreV1().Services(k.Namespace).List(context.Background(), v1.ListOptions{})

	if err != nil || services == nil {
		return
	}

	for _, svc := range services.Items {

		if strings.HasPrefix(svc.Name, constants.KubetunnelSlug) {
			continue
		}

		serviceNames = append(serviceNames, svc.Name)
	}

	return
}

func (k *Kube) PortForward(serviceName string, port string) (listeningPort int, err error) {

	service, err := k.GetServiceContext(serviceName)

	if err != nil {
		return -1, err
	}

	portInt, err := strconv.Atoi(port)
	if err != nil {
		return -1, err
	}

	pf := &portforward.PortForward{
		Namespace: k.Namespace,
		Labels: v1.LabelSelector{
			MatchLabels: service.LabelSelector,
		},
		DestinationPort: portInt,
		Config:          k.Config,
		Clientset:       k.InnerKubeClient,
	}

	listeningPort, err = pf.Start(context.Background())
	if err != nil {
		return -1, fmt.Errorf("error starting port forward: %s", err)
	}

	log.Printf("Started tunnel on %d\n", pf.ListenPort)

	// TODO: how to stop? pf.Stop()..? do we need to listen on stop??

	return listeningPort, nil
}

func (k *Kube) CreateKubeTunnelResource(resourceSpec models.KubeTunnelResourceSpec) error {
	var body []byte

	resource := models.KubeTunnelResource{
		TypeMeta: v1.TypeMeta{
			Kind:       constants.KubeTunnelKind,
			APIVersion: constants.KubeTunnelApiVersion,
		},
		Metadata: v1.ObjectMeta{
			Name: resourceSpec.KubeTunnelServiceName(),
		},
		Spec: resourceSpec,
	}

	url := fmt.Sprintf("/apis/%s/namespaces/%s/%ss", constants.KubeTunnelApiVersion, k.Namespace, constants.KubetunnelSlug)

	data, err := k.InnerKubeClient.RESTClient().
		Get().
		AbsPath(fmt.Sprintf("%s/%s", url, resource.Metadata.Name)).
		DoRaw(context.TODO())

	if err != nil {

		if body, err = json.Marshal(resource); err != nil {
			return err
		}

		data, err = k.InnerKubeClient.RESTClient().
			Post().
			AbsPath(url).
			Body(body).
			DoRaw(context.TODO())
	} else {

		var oldResource models.KubeTunnelResource
		if err = json.Unmarshal(data, &oldResource); err != nil {
			return fmt.Errorf("unable to update existing CRD. please manually delete kubetunnel: %s and run again", resource.Metadata.Name)
		}

		resource.Metadata.ResourceVersion = oldResource.Metadata.ResourceVersion

		if body, err = json.Marshal(resource); err != nil {
			return err
		}

		data, err = k.InnerKubeClient.RESTClient().
			Put().
			AbsPath(fmt.Sprintf("%s/%s", url, resource.Metadata.Name)).
			Body(body).
			DoRaw(context.TODO())
	}

	log.Debug(string(data))
	return err
}

func (k *Kube) ConnectivityCheck() error {

	_, err := k.InnerKubeClient.Discovery().ServerVersion()
	if err != nil {
		return fmt.Errorf("unable to connect kubernetes host: '%s', check KUBECONFIG set or ~/.kube/config is configured correctly", k.Config.Host)
	}

	return nil
}

func (k *Kube) RBACCheck() error {

	// Check RBAC permissions for each of the requested namespaces
	requiredPermissions := []v12.ResourceAttributes{
		{Verb: "list", Resource: "pods"},
		{Verb: "get", Resource: "pods"},
		{Verb: "watch", Resource: "pods"},
		{Verb: "get", Resource: "services"},
	}

	for _, perm := range requiredPermissions {

		var accessReview = &v12.SelfSubjectAccessReview{
			Spec: v12.SelfSubjectAccessReviewSpec{
				ResourceAttributes: &perm,
			},
		}
		accessReview, err := k.InnerKubeClient.AuthorizationV1().SelfSubjectAccessReviews().Create(context.TODO(), accessReview, v1.CreateOptions{})
		if err != nil {
			return fmt.Errorf("unable to connect kubernetes host: '%s', more info: '%s'", k.Config.Host, err.Error())
		}
		if accessReview.Status.Allowed == false {
			return fmt.Errorf("host: '%s', namespace: '%s' missing RBAC permission: %v", k.Config.Host, k.Namespace, perm)
		}
	}

	return nil
}

func createInClusterKubeClient() (*kubernetes.Clientset, *rest.Config, error) {
	inClusterConf, err := rest.InClusterConfig()

	if err != nil {
		log.Warnf("fail to get InClusterConfig, err: %s", err.Error())
		return nil, nil, err
	}

	client, err := kubernetes.NewForConfig(inClusterConf)
	if err != nil {
		log.Warnf("fail to create new config from in cluster, err: %s", err.Error())
		return nil, nil, err
	}

	log.Info("connected using in cluster config")

	return client, inClusterConf, nil
}

func mustCreateOutOfClusterKubeClient(kubeConf string) (*kubernetes.Clientset, clientcmd.ClientConfig, error) {

	var kubeConfig clientcmd.ClientConfig

	if kubeConf == "" {
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	} else {

		conf, err := clientcmd.LoadFromFile(kubeConf)

		if err != nil {
			return nil, nil, fmt.Errorf("err: unable to load kubeconfig from path: '%s'. \"%s\"", kubeConf, err.Error())
		}

		kubeConfig = clientcmd.NewDefaultClientConfig(*conf, &clientcmd.ConfigOverrides{})
	}

	config, err := kubeConfig.ClientConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("err: unable to read kubeconfig. \"%s\"", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, nil, fmt.Errorf("err: unable to create kube client. \"%s\"", err.Error())
	}

	return kubeClient, kubeConfig, nil
}
