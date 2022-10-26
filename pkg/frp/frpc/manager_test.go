package frpc_test

import (
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"testing"
)

func TestFrpClientManager(t *testing.T) {

	common := models.Common{
		ServerAddress: "kubetunnel-demo",
		ServerPort:    "7000",
	}

	svc := []frpc.ServicePair{
		{
			Name: "local_app",
			Service: models.Service{
				Type:       "tcp",
				RemotePort: "80",
				LocalIP:    "localhost",
				LocalPort:  "22285",
			},
		},
	}

	manager := frpc.NewManager(common, svc, nil)

	manager.RunFRPc()

}
