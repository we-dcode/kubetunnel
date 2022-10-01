package pkg

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/txn2/kubefwd/pkg/utils"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube/servicecontext"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	frpmodels "github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/kubefwd"
)

type KubeTunnel struct {
	kubeClient *kube.Kube
	helmClient *helm.Helm
}

type KubeTunnelConf struct {
	GCVersion            string
	FRPSVersion          string
	ServiceName          string
	RemoteToLocalPortMap map[string]string
	LocalIP              string
}

func MustNew(namespace string) *KubeTunnel {

	kubeClient := kube.MustNew(namespace)

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

	return &KubeTunnel{
		kubeClient: kubeClient,
		helmClient: helmClient,
	}
}

func (ct *KubeTunnel) Run(tunnelConf KubeTunnelConf) {

	err := ct.helmClient.InstallOrUpgradeGC(tunnelConf.GCVersion)
	if err != nil {
		log.Panic(err.Error())
	}

	serviceContext, err := ct.kubeClient.GetServiceContext(tunnelConf.ServiceName)
	if err != nil {
		log.Panic(err.Error())
	}

	frpServerValues := servicecontext.ToFRPServerValues(serviceContext)

	err = ct.helmClient.InstallOrUpgradeFrpServer(tunnelConf.FRPSVersion, frpServerValues)
	if err != nil {
		log.Panic(err.Error())
	}

	// TODO: The following function is blocking. Need to create a channel and re-use it for both kubefwd and frpc
	err = kubefwd.Execute(ct.kubeClient)
	if err != nil {
		log.Panicf("fail executing kubefwd: %s", err.Error())
	}

	common := frpmodels.Common{
		ServerAddress: fmt.Sprintf("%s-frps", frpServerValues.ServiceName),
		ServerPort:    constants.FRPServerPort,
	}

	servicePortsPairs := servicecontext.ToFRPClientPairs(tunnelConf.LocalIP, tunnelConf.RemoteToLocalPortMap, serviceContext)

	err = frpc.Execute(common, servicePortsPairs...)

	if err != nil {
		log.Panic(err.Error())
	}
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
