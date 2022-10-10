package frpc

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"github.com/txn2/kubefwd/pkg/fwdport"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/frp/models"
	"github.com/we-dcode/kube-tunnel/pkg/notify"
	"github.com/we-dcode/kube-tunnel/pkg/notify/killsignal"
	"github.com/we-dcode/kube-tunnel/pkg/utils/hostsutils"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	"strings"
	"time"
)

type Manager struct {
	Common      models.Common
	ServicePair []ServicePair
	HostFile    *fwdport.HostFileWithLock
}

func NewManager(common models.Common, servicePair []ServicePair, hostFile *fwdport.HostFileWithLock) *Manager {

	return &Manager{
		common,
		servicePair,
		hostFile,
	}
}

func (m *Manager) RunFRPc() {

	for killsignal.HasKillSignaled() == false {

		m.WaitForLocalPortToBecomeAvailable()

		cancelChan, err := Execute(m.Common, m.ServicePair)
		if err != nil {
			log.Panicf("unable to start frp client, err: %s", err.Error())
		}

		// Known limitation: host will never be rolled back to original dns record.
		// We expect the user to use kubetunnel only when he wants to forward traffic to his local env.
		//ChangeHostToKubeTunnel(m.HostFile, m.Common.ServerAddress)

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

func ChangeHostToKubeTunnel(hostFile *fwdport.HostFileWithLock, kubeTunnelServerDns string) {

	originalServiceDns := strings.Replace(kubeTunnelServerDns, fmt.Sprintf("%s-", constants.KubetunnelSlug), "", 1)

	hostsutils.ReplaceAddressForHost(hostFile.Hosts, originalServiceDns, kubeTunnelServerDns)
}
