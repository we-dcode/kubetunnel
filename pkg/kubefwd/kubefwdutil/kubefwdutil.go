package kubefwdutil

import (
	log "github.com/sirupsen/logrus"
	"github.com/txn2/txeh"
	"strings"
)

func HostsCleanup(hosts *txeh.Hosts)  {

	log.Info("cleaning up all entries containing .kubetunnel host")

	for _, line := range *hosts.GetHostFileLines() {

		for _, hostname := range line.Hostnames {

			if strings.Contains(hostname, "kubetunnel") {

				hosts.RemoveAddress(line.Address)
				break
			}
		}
	}

}