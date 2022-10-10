package kubefwdutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"strings"
)

func HostsCleanup(hosts *txeh.Hosts) {

	log.Info("cleaning up all entries containing .kubetunnel host")

	for _, line := range *hosts.GetHostFileLines() {

		for _, hostname := range line.Hostnames {

			if strings.Contains(hostname, constants.KubetunnelSlug) {

				hosts.RemoveAddress(line.Address)
				break
			}
		}
	}

}
