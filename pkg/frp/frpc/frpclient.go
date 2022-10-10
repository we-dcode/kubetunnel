package frpc

import (
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frputil"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/notify"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tomlutil"
)

type ServicePair struct {
	Name    string
	Service models.Service
}

// Execute - This code was copied from frpc and modified a bit to support kubetunnel requirements
func Execute(common models.Common, servicePair []ServicePair) (cancelChan *notify.CancellationChannel, err error) {

	tomlString := createToml(common, servicePair)

	cfg, pxyCfgs, visitorCfgs, err := frputil.ParseClientConfig([]byte(tomlString))
	if err != nil {

		return nil, fmt.Errorf("fail to start frpc. more info: '%s'", err.Error())
	}

	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, "")
	if errRet != nil {
		err = errRet
		return nil, err
	}

	cancelChan = notify.NewCancellationChannelWithCallback(func() {
		svr.Close()
	})

	go func() {
		err = svr.Run()
		cancelChan.CancelWithReason(err) // if finish with reason
	}()

	return cancelChan, nil
}

func createToml(common models.Common, servicePair []ServicePair) string {
	frpConfig := models.FrpClientConfig{
		"common": common,
	}

	for _, element := range servicePair {

		frpConfig[element.Name] = element.Service
	}

	tomlString, _ := tomlutil.Marshal(frpConfig)

	return tomlString
}
