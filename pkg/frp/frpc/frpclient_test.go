package frpc_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"testing"
)

func TestInstallingKubetunnelGC(t *testing.T) {

	common := models.Common{
		ServerAddress: "localhost",
		ServerPort:    "7001",
	}

	svc := []frpc.ServicePair{
			{
				Name: "google",
				Service: models.Service{
					Type:       "tcp",
					RemotePort: "80",
					LocalIP:    "google.com",
					LocalPort:  "80",
				},
			},
			{
				Name: "microsoft",
				Service: models.Service{
				Type:       "tcp",
				RemotePort: "8081",
				LocalIP:    "microsoft.com",
				LocalPort:  "80",
			},
		},
	}

	err := frpc.Execute(common, svc...)

	assert.NoError(t, err)
}