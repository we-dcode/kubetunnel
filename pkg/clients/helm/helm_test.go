package helm_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	models2 "github.com/we-dcode/kube-tunnel/pkg/clients/helm/models"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"testing"
)

func TestInstallingKubetunnelFrp(t *testing.T) {

	client := kube.MustNew("")

	helmClient := helm.MustNew(client)

	err := helmClient.InstallOrUpgradeFrpServer("", &models2.FRPServerValues{
		Ports:             models2.Ports{
			Values: []string{ "8081", "8082" },
		},
		ServiceName:       "kubetunnel-svc",
		PodSelectorLabels: []models2.PodSelectorLabel{
			{Key: "app", Value: "kubetunnel-svc"},
		},
	})

	assert.NoError(t, err)
}


func TestInstallingKubetunnelGC(t *testing.T) {

	client := kube.MustNew("")

	helmClient := helm.MustNew(client)

	err := helmClient.InstallOrUpgradeGC("")

	assert.NoError(t, err)
}