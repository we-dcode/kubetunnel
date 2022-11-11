package pkg_test

import (
	"github.com/we-dcode/kube-tunnel/pkg"
	"testing"
)

func TestRunningKubeTunnelE2E2(t *testing.T) {

	kubeTunnel := pkg.MustNewKubeTunnel("/Users/maordavidov/dcode/gitlab-cicd-kubeconfig.yaml", "kubetunnel", true)

	//kubeTunnel.Install("0.0.12")

	kubeTunnel.CreateTunnel(pkg.KubeTunnelConf{
		ServiceName: "nginx",
		LocalIP:     "localhost",
		KubeTunnelPortMap: map[string]string{
			"8081": "80",
		},
		DnsForwardAllNamespaces: true,
		Proxies: map[string]int{
			"postgres.dfv7txrpzu6m.us-east-1.rds.amazonaws.com": 5432,
		},
	})

	//assert.NoError(t, err)
}
