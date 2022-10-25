package kube_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGetServiceWithSinglePortOnDefaultNamespace(t *testing.T) {

	client := kube.MustNew("", "")

	context, err := client.GetServiceContext("kubetunnel-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 1)
}

func TestGetServiceWithMultiplePortsOnDefaultNamespace(t *testing.T) {

	client := kube.MustNew("", "")

	context, err := client.GetServiceContext("kubetunnel-multi-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 2)
}

func TestGetServiceFromExplicitNamespaceWithMultipleLables(t *testing.T) {

	client := kube.MustNew("", "kubetunnel-explicit")

	context, err := client.GetServiceContext("kubetunnel-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 1)
	assert.Len(t, context.LabelSelector, 2)
}

func TestKubePortForward(t *testing.T) {

	client := kube.MustNew("/Users/maordavidov/dcode/gitlab-cicd-kubeconfig.yaml", "kubetunnel")

	listeningPort, err := client.PortForward("kubetunnel-nginx", "7000")

	assert.NoError(t, err)
	assert.NotEqual(t, -1, listeningPort)
}

func TestCreateKubeTunnelResource(t *testing.T) {

	client := kube.MustNew("/Users/maordavidov/dcode/gitlab-cicd-kubeconfig.yaml", "kubetunnel")

	err := client.CreateKubeTunnelResource(kube.KubeTunnelResource{
		TypeMeta: metav1.TypeMeta{
			Kind:       constants.KubetunnelSlug,
			APIVersion: constants.KubeTunnelApiVersion,
		},
		ObjectMeta: metav1.ObjectMeta{
			Name: "kubetunnel-nginx",
		},
		Spec: kube.KubeTunnelResourceSpec{
			Ports:       kube.Ports{Values: []string{"80"}},
			ServiceName: "nginx",
			PodSelectorLabels: map[string]string{
				"app": "nginx",
			},
		},
	})

	assert.NoError(t, err)
}
