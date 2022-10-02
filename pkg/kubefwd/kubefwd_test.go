package kubefwd_test

import (
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm/models"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/kubefwd"
	"testing"
)

// TODO: execute this test using sudo -E go test -run 'KubeFwd'. OR configure Goland to run as sudo
func TestKubeFwd(t *testing.T) {

	kubeClient := kube.MustNew("kubetunnel")

	// kubectl run nginx -l app=nginx --image=nginx:1.14.2 --port 80
	// kubectl run client-mock -l app=client-mock --image=nginx:1.14.2
	// kubectl create svc clusterip nginx --tcp 80:80
	// kubectl create svc clusterip kubetunnel-nginx --tcp 80:80 // TODO: change label selector to app: nginx
	frpsValues := &models.FRPServerValues{
		Ports:             models.Ports{
			Values: []string{ "80" },
		},
		ServiceName:       "nginx",
	}

	errChan := kubefwd.Execute(kubeClient, frpsValues)

	log.Info("waiting for channel to complete...")

	assert.NoError(t, <-errChan)
}