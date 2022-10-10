package frpc

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/notify"
	"github.com/we-dcode/kube-tunnel/pkg/notify/killsignal"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	"time"
)

type Manager struct {
	Common      models.Common
	ServicePair []ServicePair
}

func NewManager(common models.Common, servicePair []ServicePair) *Manager {

	//hostFile *fwdport.HostFileWithLock

	return &Manager{
		common,
		servicePair,
	}
}

func (m *Manager) RunFRPc() {

	for killsignal.HasKillSignaled() == false {

		m.WaitForLocalPortToBecomeAvailable()

		cancelChan, err := Execute(m.Common, m.ServicePair)
		if err != nil {
			log.Panicf("unable to start frp client, err: %s", err.Error())
		}

		ModifyHosts()

		go func() {
			m.WaitForLocalPortToBecomeUnavailableAndCancel(cancelChan)
		}()

		cancelChan.WaitForCancellation()
	}

}

func (m *Manager) WaitForLocalPortToBecomeAvailable() {

	// TODO: assuming we are using a single service at the time
	host := m.ServicePair[0].Service.LocalIP
	port := m.ServicePair[0].Service.LocalPort

	for killsignal.HasKillSignaled() == false && tcputil.IsAvailable(host, port) == false {
		time.Sleep(time.Millisecond * 500)
	}
}

func (m *Manager) WaitForLocalPortToBecomeUnavailableAndCancel(channel *notify.CancellationChannel) {

	// TODO: assuming we are using a single service at the time
	host := m.ServicePair[0].Service.LocalIP
	port := m.ServicePair[0].Service.LocalPort

	for killsignal.HasKillSignaled() == false && channel.IsCancelled() == false && tcputil.IsAvailable(host, port) {
		time.Sleep(time.Millisecond * 500)
	}

	if killsignal.HasKillSignaled() == false && channel.IsCancelled() == false {
		channel.CancelWithReason(fmt.Errorf("service is unavailable at address: %s:%s, shutting down frpc", host, port))
	}
}

func ModifyHosts() {

}
