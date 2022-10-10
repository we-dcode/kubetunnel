package frpc

import (
	"fmt"
	"github.com/fatedier/frp/client"
	"github.com/txn2/kubefwd/pkg/fwdport"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frputil"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tomlutil"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type ServicePair struct {
	Name    string
	Service models.Service
}

// Execute - This code was copied from frpc and modified a bit to support kubetunnel requirements
func Execute(common models.Common, hostFile *fwdport.HostFileWithLock, servicePair ...ServicePair) (err error) {

	frpConfig := models.FrpClientConfig{
		"common": common,
	}

	for _, element := range servicePair {

		frpConfig[element.Name] = element.Service
	}

	tomlString, err := tomlutil.Marshal(frpConfig)

	cfg, pxyCfgs, visitorCfgs, err := frputil.ParseClientConfig([]byte(tomlString))
	if err != nil {

		return fmt.Errorf("fail to start frpc. more info: '%s'", err.Error())
	}

	svr, errRet := client.NewService(cfg, pxyCfgs, visitorCfgs, "")
	if errRet != nil {
		err = errRet
		return err
	}

	kcpDoneCh := make(chan struct{})
	// Capture the exit signal if we use kcp.
	if cfg.Protocol == "kcp" {
		go handleSignal(svr, kcpDoneCh)
	}

	err = svr.Run()
	if err == nil && cfg.Protocol == "kcp" {
		<-kcpDoneCh
	}

	return nil
}
func handleSignal(svr *client.Service, doneCh chan struct{}) {
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	<-ch
	svr.GracefulClose(500 * time.Millisecond)
	close(doneCh)
}
