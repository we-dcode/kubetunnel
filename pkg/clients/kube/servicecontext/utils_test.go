package servicecontext_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube/servicecontext"
	v1 "k8s.io/api/core/v1"
	"testing"
)

func TestServiceContextToFRPSValues(t *testing.T) {

	svcContext := servicecontext.ServiceContext{
		ServiceName: "kubetunnel-svc-test",
		LabelSelector: map[string]string{
			"kuku": "riku",
			"abra": "kadabra",
			"foo":  "bar",
		},
		Ports: []v1.ServicePort{
			{
				Port: 80,
			},
			{
				Port: 443,
			},
		},
	}

	values := servicecontext.ToFRPServerValues(&svcContext)

	assert.NotNil(t, values)

	assert.Equal(t, svcContext.ServiceName, values.ServiceName)
	assert.Len(t, values.Ports.Values, 2)
	assert.Len(t, values.PodSelectorLabels, 3)
}

func TestWhenCallingToFRPClientPairs_AndNotAllPortMapFound_MapOnlyExistingPorts(t *testing.T) {

	svcContext := servicecontext.ServiceContext{
		ServiceName: "kubetunnel-svc-test",
		Ports: []v1.ServicePort{
			{
				Port: 80,
			},
			{
				Port: 443,
			},
		},
	}

	portMapping := map[string]string{
		"56172": "80",
	}

	values := servicecontext.ToFRPClientPairs("127.0.0.1", portMapping, &svcContext)

	assert.NotNil(t, values)
	assert.Len(t, values, 1)
}

func TestWhenCallingToFRPClientPairs_AndAllPortMapFound_MapAllPorts(t *testing.T) {

	svcContext := servicecontext.ServiceContext{
		ServiceName: "kubetunnel-svc-test",
		Ports: []v1.ServicePort{
			{
				Port: 80,
			},
			{
				Port: 443,
			},
		},
	}

	portMapping := map[string]string{
		"56172": "80",
		"56171": "443",
	}

	values := servicecontext.ToFRPClientPairs("127.0.0.1", portMapping, &svcContext)

	assert.NotNil(t, values)
	assert.Len(t, values, 2)
}
