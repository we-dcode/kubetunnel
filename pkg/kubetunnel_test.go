package pkg_test

import (
	"github.com/we-dcode/kube-tunnel/pkg"
	"testing"
)

func TestRunningKubeTunnelE2E(t *testing.T) {

	kubeTunnel := pkg.MustNewKubeTunnel("kubetunnel")

	kubeTunnel.Run(pkg.KubeTunnelConf{
		GCVersion:         "0.1.1",
		KubeTunnelVersion: "0.1.2",
		ServiceName:       "nginx",
		LocalIP:           "localhost",
		KubeTunnelPortMap: map[string]string{
			"8081": "80",
		},
	})

	//assert.NoError(t, err)
}
