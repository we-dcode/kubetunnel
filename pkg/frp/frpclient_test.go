package frp_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"testing"
)

func TestInstallingKubetunnelGC(t *testing.T) {

	client := kube.MustNew("")

	helmClient := helm.MustNew(client)

	err := helmClient.InstallOrUpgradeGC("")

	assert.NoError(t, err)
}