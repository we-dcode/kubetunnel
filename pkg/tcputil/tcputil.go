package tcputil

import (
	log "github.com/sirupsen/logrus"
	"net"
	"time"
)

func IsAvailable(host string, port string) bool{

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	if err != nil {
		log.Debugf("connection error: %s", err)
		return false
	}

	defer conn.Close()
	log.Debugf("connection succeeded: %s", net.JoinHostPort(host, port))

	return true
}