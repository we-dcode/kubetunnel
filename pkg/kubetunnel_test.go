package pkg_test

import (
	"github.com/we-dcode/kube-tunnel/pkg"
	"testing"
)

func TestRunningKubeTunnelE2E2(t *testing.T) {

	kubeTunnel := pkg.MustNewKubeTunnel("/Users/maordavidov/dcode/gitlab-cicd-kubeconfig.yaml", "kubetunnel")

	kubeTunnel.CreateTunnel(pkg.KubeTunnelConf{
		ServiceName: "nginx",
		LocalIP:     "localhost",
		KubeTunnelPortMap: map[string]string{
			"8081": "80",
		},
	})

	//assert.NoError(t, err)
}
