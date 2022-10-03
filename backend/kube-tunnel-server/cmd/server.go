package main

import (
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// album represents data about a record album.
type Environment struct {
	ServiceName       string            `yaml:"env_service_name,omitempty"`
	PodSelectorLabels map[string]string `yaml:"pod_selector_labels,omitempty"`
}

var ports = []int{8080, 8081}

func main() {
	router := gin.Default()
	router.GET("/health", healthHandler)

	router.Run("localhost:8080")
}

func healthHandler(c *gin.Context) {

	for _, port := range ports {
		if !isAvailable("127.0.0.1", strconv.FormatInt(int64(port), 10)) {
			log.Debugf("port %v is unavailable on host", port)
			c.IndentedJSON(http.StatusOK, "Fail")
		}
	}
	c.IndentedJSON(http.StatusOK, ports)
}

func isAvailable(host string, port string) bool {

	timeout := time.Second
	conn, err := net.DialTimeout("tcp", net.JoinHostPort(host, port), timeout)
	log.Debug("bla")
	if err != nil {
		log.Debugf("connection error: %s", err)
		return false
	}

	defer conn.Close()
	log.Debugf("connection succeeded: %s", net.JoinHostPort(host, port))

	return true
}
