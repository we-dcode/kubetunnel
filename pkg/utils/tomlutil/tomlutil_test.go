package tomlutil_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tomlutil"
	"testing"
)

func TestMarshalTomlWithCommonOnly(t *testing.T) {

	frpConfig := models.FrpClientConfig{
		"Common": models.Common{
			ServerAddress: "localhost",
			ServerPort:    "7001",
		},
	}

	tomlString, err := tomlutil.Marshal(frpConfig)

	assert.NoError(t, err)
	assert.NotEmpty(t, tomlString)
}

func TestMarshalTomlWithMultipleServices(t *testing.T) {

	frpConfig := models.FrpClientConfig{
		"common": models.Common{
			ServerAddress: "localhost",
			ServerPort:    "7001",
		},
		"rabbit": models.Service{
			Type:       "tcp",
			RemotePort: "5672",
			LocalIP:    "localhost",
			LocalPort:  "5672",
		},
		"rabbit_management_ui": models.Service{
			Type:       "tcp",
			RemotePort: "15672",
			LocalIP:    "localhost",
			LocalPort:  "15672",
		},
	}

	tomlString, err := tomlutil.Marshal(frpConfig)

	assert.NoError(t, err)
	assert.NotEmpty(t, tomlString)
}
