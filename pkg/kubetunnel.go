package pkg

import (
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
)

func KubeTunnelExecute()  {

	kubeClient := kube.MustNew("")

	err := kubeClient.ConnectivityCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	err = kubeClient.RBACCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	helmClient := helm.MustNew(kubeClient)




	_ = helmClient
}
