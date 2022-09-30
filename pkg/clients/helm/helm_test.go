package helm_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"testing"
)

func TestInstallingOperator(t *testing.T) {

	client := kube.MustNew("/Users/maordavidov/dcode/gitlab-cicd-kubeconfig.yaml", "kubetunnel")

	helmClient := helm.MustNew(client)

	err := helmClient.InstallKubeTunnelOperator("0.0.6")

	assert.NoError(t, err)
}
