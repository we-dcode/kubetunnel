package kube

import (
	"context"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {

	kubeConfig clientcmd.ClientConfig
	kubeClient kubernetes.Interface
	config *rest.Config
}


func MustNew() *KubeClient {

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	configOverrides := &clientcmd.ConfigOverrides{}
	kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
	config, err := kubeConfig.ClientConfig()
	if err != nil {
		log.Panicf("err: unable to read kubeconfig. \"%s\"", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if  err != nil {
		log.Panicf("err: unable to create kube client. \"%s\"", err.Error())
	}

	return &KubeClient{
		kubeConfig,
		kubeClient,
		config,
	}
}

func (k *KubeClient) GetServiceContext(name string, namespace string) ServiceContext  {

	if len(namespace) == 0 {
		namespace, _, _ = k.kubeConfig.Namespace()
	}

	svc, err := k.kubeClient.CoreV1().Services(namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {
		log.Panicf("namespace: '%s' svc: '%s' not found at host: '%s'", namespace, name, k.config.Host)
	}

	return ServiceContext{
		LabelSelector: svc.Spec.Selector,
		Ports:         svc.Spec.Ports,
	}
}

