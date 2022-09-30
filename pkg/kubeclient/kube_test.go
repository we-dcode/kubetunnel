package kubeclient_test

import (
	"github.com/we-dcode/kube-tunnel/pkg/kubeclient"
	"github.com/we-dcode/kube-tunnel/pkg/models"
	"testing"
)

func TestInstallHelm(t *testing.T) {

 	client := kubeclient.MustNew()

	client.MustInstallFrpServer("", &models.FrpServerValues{
		Namespace:         "",
		Ports:             models.Ports{
			Values: nil,
		},
		ServiceName:       "",
		PodSelectorLabels: nil,
	})

	//assert.Equal(t, count, atomic.LoadInt32(&count))
}

