package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/robfig/cron/v3"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/utils/logutil"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"net/http"
	"os"
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
	//s := gocron.NewScheduler(time.UTC)

	//s.Every(5).Seconds().Do(func() { portChecker() })
	job := cron.New()
	job.AddFunc("@every 5s", func() {
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
	kube := connectToKubernetes() // TODO: change this one

	portArr := strings.Split(getEnvVar("PORTS"), ",")
	serviceName := getEnvVar("SERVICE_NAME")

	log.Debugf("ports are %v", portArr)

	for _, port := range portArr {
		if tcputil.IsAvailable(hostname, port) == false {
			log.Debugf("port %v is unavailable on host", port)

			// TODO Make the operator do the patching of the service
			error := patchServiceWithLabel(kube, serviceName, false)
			if error != nil {
				log.Debugf("error %v", error)
			} else {
				isConnected = false
			}
			return
		}
	}

	//  patch
	log.Debugf("labeling service %v to %v\n", serviceName, true)
	error := patchServiceWithLabel(kube, serviceName, true)
	if error != nil {
		log.Debugf("error %v", error)
	} else {
		isConnected = true
	}

}

func healthHandler(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, isConnected)
}

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

func getEnvVar(variable string) string {
	envVar, ok := os.LookupEnv(variable)
	if !ok {
		log.Errorf("%v is not a present env variable", variable)
	}
	return envVar
}

func connectToKubernetes() *kube.Kube {

	kubeClient := kube.MustNew("", "kubetunnel")

	err := kubeClient.ConnectivityCheck()
	if err != nil {
		log.Errorf(err.Error())
	}

	//err = kubeClient.RBACCheck()
	//if err != nil {
	//	log.Errorf(err.Error())
	//}

	return kubeClient
}

func patchServiceWithLabel(kube *kube.Kube, serviceName string, connected bool) error {

	svcContext, err := kube.GetServiceContext(serviceName)
	if err != nil {
		log.Errorf("fail to get service: '%s' context, error: %s", serviceName, err.Error())
	}

	clientSet := kube.InnerKubeClient
	log.Debugf(kube.Namespace)
	ctx := context.TODO()

	slugPrefix := fmt.Sprintf("%s-", constants.KubetunnelSlug)

	// TODO: Check if service needs to be updated...
	//serviceUpdateRequires := false

	if !connected {

		log.Debugf("removing true from %v\n", serviceName)
		var payload []patchStringValue
		//payload := []patchStringValue{{ TODO: This label should be added to service label and not pod label selector
		//	Op:    "remove",
		//	Path:  "/spec/selector/kube-tunnel",
		//	Value: "true",
		//}}

		for key, valueWithSlug := range svcContext.LabelSelector {

			if strings.EqualFold(key, constants.KubetunnelSlug) {
				continue
			}

			slugAlreadyRemoved := strings.Contains(valueWithSlug, constants.KubetunnelSlug) == false
			if slugAlreadyRemoved {
				continue
			}

			valueWithoutSlug := strings.Replace(valueWithSlug, slugPrefix, "", 1)

			payload = append(payload, patchStringValue{

				Op:    "replace",
				Path:  fmt.Sprintf("/spec/selector/%s", key),
				Value: valueWithoutSlug,
			})
		}

		payloadBytes, _ := json.Marshal(payload)
		_, err := clientSet.
			CoreV1().
			Services(kube.Namespace).
			Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
		return err
	} else {
		log.Debugf("adding true to %v\n", serviceName)
		//payload := []patchStringValue{{
		//	Op:    "add",
		//	Path:  "/spec/selector/kube-tunnel",
		//	Value: "true",
		//}}

		var payload []patchStringValue

		for key, valueWithoutSlug := range svcContext.LabelSelector {

			if strings.EqualFold(key, constants.KubetunnelSlug) {
				continue
			}

			alreadyContainSlug := strings.Contains(valueWithoutSlug, constants.KubetunnelSlug)

			if alreadyContainSlug {
				continue
			}

			valueWithSlug := fmt.Sprintf("%s%s", slugPrefix, valueWithoutSlug)

			payload = append(payload, patchStringValue{

				Op:    "replace",
				Path:  fmt.Sprintf("/spec/selector/%s", key),
				Value: valueWithSlug,
			})
		}

		payloadBytes, _ := json.Marshal(payload)
		_, err := clientSet.
			CoreV1().
			Services(kube.Namespace).
			Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
		return err
	}
}
