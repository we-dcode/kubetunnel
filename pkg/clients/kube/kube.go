package kube

import (
	"context"
	"fmt"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type Kube struct {

	kubeConfig clientcmd.ClientConfig
	kubeClient kubernetes.Interface
	Config *rest.Config
	Namespace string
}


func MustNew(namespace string) *Kube {

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

	if len(namespace) == 0 {
		namespace, _, _ = kubeConfig.Namespace()
	}

	return &Kube{
		kubeConfig,
		kubeClient,
		config,
		namespace,
	}
}

func (k *Kube) GetServiceContext(name string) (*ServiceContext, error)  {

	svc, err := k.kubeClient.CoreV1().Services(k.Namespace).Get(context.Background(), name, v1.GetOptions{})
	if err != nil {

		err = fmt.Errorf("namespace: '%s' svc: '%s' not found at host: '%s'", k.Namespace, name, k.Config.Host)
		return nil, err
	}

	ctx := ServiceContext{
		LabelSelector: svc.Spec.Selector,
		Ports:         svc.Spec.Ports,
	}

	return &ctx, nil
}

