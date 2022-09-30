package frputil_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frputil"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/tomlutil"
	"testing"
)

func TestGetServiceWithSinglePortOnDefaultNamespace(t *testing.T) {

	frpConfig := models.FrpClientConfig{
		"common":  models.Common{
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

	conf, pxyCfgs, visitorCfgs, err := frputil.ParseClientConfig([]byte(tomlString))

	assert.NotNil(t, conf)
	assert.NotNil(t, pxyCfgs)
	assert.NotNil(t, visitorCfgs)
	assert.NoError(t, err)
}