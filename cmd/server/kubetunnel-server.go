package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/utils/logutil"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
	"os"
	"strconv"
	"strings"
)

var (
	kubeContext = ""

	isConnected = false

	hostname = "127.0.0.1"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})
	log.SetOutput(&logutil.LogOutputSplitter{})
	log.Print("")
	log.Print("https://github.com/we-dcode/kube-tunnel")
	log.Print("https://we.dcode.tech")
	log.Print("")

	startJob()
	startGin()

}

func startJob() {
	job := cron.New()
	job.AddFunc("@every 1s", func() {
		portChecker()
	})
	job.Start()
}

func startGin() {
	router := gin.Default()
	router.GET("/health", healthHandler)

	router.Run("0.0.0.0:8080")
}

func portChecker() {

	log.Debugf("Starting portChecker.. ")

	portArr := strings.Split(getEnvVar("PORTS"), ",")
	serviceName := getEnvVar("SERVICE_NAME")
	podNamespace := getEnvVar("POD_NAMESPACE")
	operatorSvcName := getEnvVar("OPERATOR_SVC_NAME")
	operatorNamespace := getEnvVar("OPERATOR_NAMESPACE")
	operatorPort := getEnvVar("OPERATOR_PORT")

	log.Debugf("ports are %v", portArr)

	for _, port := range portArr {
		if tcputil.IsAvailable(hostname, port) == false {
			log.Debugf("port %v is unavailable on host", port)

			if isConnected != false {
				error := patchService(
					podNamespace, serviceName, false, operatorSvcName,
					operatorNamespace, operatorPort)
				if error != nil {
					log.Debugf("error %v", error)
				} else {
					isConnected = false
				}

			}

			return
		}
	}

	//  patch
	log.Debugf("labeling service %v to %v\n", serviceName, true)
	if isConnected != true {
		error := patchService(
			podNamespace, serviceName,
			true, operatorSvcName,
			operatorNamespace, operatorPort)
		if error != nil {
			log.Debugf("error %v", error)
		} else {
			isConnected = true
		}
	}

}

func healthHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, isConnected)
}

func patchService(namespace string, serviceName string, isConnected bool, operatorSvcName string, operatorNamespace string, operatorSvcPort string) error {
	values := map[string]string{"Namespace": namespace,
		"Name":        serviceName,
		"IsConnected": strconv.FormatBool(isConnected)}
	jsonData, err := json.Marshal(values)

	if err != nil {
		return err
	}

	resp, err := http.Post(fmt.Sprintf(
		"http://%s.%s.svc.cluster.local:%s/service", operatorSvcName, operatorNamespace, operatorSvcPort),
		"application/json",
		bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}
	log.Printf("Response: %s", resp)
	return nil

}
func getEnvVar(variable string) string {
	envVar, ok := os.LookupEnv(variable)
	if !ok {
		log.Errorf("%v is not a present env variable", variable)
	}
	return envVar
}
