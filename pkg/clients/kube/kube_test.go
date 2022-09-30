package kube_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"testing"
)

func TestGetServiceWithSinglePortOnDefaultNamespace(t *testing.T) {

	client := kube.MustNew("")

	context, err := client.GetServiceContext("kubetunnel-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 1)
}

func TestGetServiceWithMultiplePortsOnDefaultNamespace(t *testing.T) {

	client := kube.MustNew("")

	context, err := client.GetServiceContext("kubetunnel-multi-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 2)
}


func TestGetServiceFromExplicitNamespaceWithMultipleLables(t *testing.T) {

	client := kube.MustNew("kubetunnel-explicit")

	context, err := client.GetServiceContext("kubetunnel-svc")

	assert.NoError(t, err)
	assert.NotNil(t, context)
	assert.Len(t, context.Ports, 1)
	assert.Len(t, context.LabelSelector, 2)
}
