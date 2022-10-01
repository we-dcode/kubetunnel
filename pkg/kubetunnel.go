package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/txn2/kubefwd/pkg/utils"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
)

func KubeTunnelExecute() {

	kubeClient := kube.MustNew("")

	err := kubeClient.ConnectivityCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	err = kubeClient.RBACCheck()
	if err != nil {
		log.Panicf(err.Error())
	}

	err = CheckRootPermissions()
	if err != nil {
		log.Panicf(err.Error())
	}

	helmClient := helm.MustNew(kubeClient)

	_ = helmClient
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
