package killsignal

import (
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/notify"
	"os"
	"os/signal"
	"syscall"
)

var (
	CancellationChannel *notify.CancellationChannel
)

func init() {

	CancellationChannel = notify.NewCancellationChannel()

	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGTERM)
		defer func() {
			signal.Stop(sigint)
		}()
		<-sigint
		log.Infof("Received shutdown signal")
		CancellationChannel.Cancel()
	}()
}

func WaitForKillSignal() {
	CancellationChannel.WaitForCancellation()
}

func HasKillSignaled() bool {
	return CancellationChannel.IsCancelled()
}
