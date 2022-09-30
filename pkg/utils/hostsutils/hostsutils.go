package hostsutils

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

func ReplaceAddressForHost(hosts *txeh.Hosts, srcHost, dstHost string) {

	log.Infof("replacing host: '%s' with: '%s'", srcHost, dstHost)

	found, newAddr, _ := hosts.HostAddressLookup(dstHost)

	if found == false {
		log.Panicf("unable to locate host: %s in hosts file. please run %s again.", dstHost, constants.KubetunnelSlug)
	}

	for _, line := range *hosts.GetHostFileLines() {

		for _, hostname := range line.Hostnames {

			if strings.EqualFold(srcHost, hostname) {

				hosts.RemoveAddress(line.Address)
				hosts.AddHosts(newAddr, line.Hostnames)

				err := hosts.Save()
				if err != nil {
					log.Panicf("unable to save host: %s in hosts file. please run %s again. internal err: %s", dstHost, constants.KubetunnelSlug, err.Error())
				}

				break
			}
		}
	}
}
