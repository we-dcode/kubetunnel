package main

import (
	"bytes"
	"context"

	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	"k8s.io/client-go/tools/clientcmd"
	"net"
	"net/http"
	"os"
	"strings"
	"time"
)

var (
	kubeContext = ""

	hostname = "127.0.0.1"
)

func main() {

	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		ForceColors:     true,
		TimestampFormat: "15:04:05",
	})
	log.SetOutput(&LogOutputSplitter{})
	log.Print("")
	log.Print("https://github.com/we-dcode/kube-tunnel")
	log.Print("https://we.dcode.tech")
	log.Print("")

	router := gin.Default()
	router.GET("/health", healthHandler)

	router.Run("localhost:8080")
}

func healthHandler(c *gin.Context) {

	clientSet := connectToKubernetes()
	portArr := strings.Split(getEnvVar("PORTS"), ",")
	serviceName := getEnvVar("SERVICE_NAME")

	log.Debugf("ports are %v", portArr)

	for _, port := range portArr {
		if !isPortAvailable(hostname, port) {
			log.Debugf("port %v is unavailable on host", port)
			//  Scale our replication controller.
			fmt.Printf("labeling service %v to connected: %v\n", serviceName, false)
			error := patchServiceWithLabel(clientSet, serviceName, false)
			if error != nil {
				log.Panicf("error %v", error)
			}
			c.IndentedJSON(http.StatusOK, "Fail")
			return
		}
	}

	//  patch
	log.Debugf("labeling service %v to %v\n", serviceName, true)
	error := patchServiceWithLabel(clientSet, serviceName, true)
	if error != nil {
		log.Panicf("error %v", error)
	}
	c.IndentedJSON(http.StatusOK, portArr)
}

//  patchStringValue specifies a patch operation for a string.
type patchStringValue struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value string `json:"value"`
}

//  patchStringValue specifies a patch operation for a uint32.
type patchUInt32Value struct {
	Op    string `json:"op"`
	Path  string `json:"path"`
	Value uint32 `json:"value"`
}

func getEnvVar(variable string) string {
	envVar, ok := os.LookupEnv(variable)
	if !ok {
		log.Panicf("%v is not a present env variable", variable)
	}
	return envVar
}

func isPortAvailable(host string, port string) bool {

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

func connectToKubernetes() *kubernetes.Clientset {

	//  Get the local kube config.
	fmt.Printf("Connecting to Kubernetes Context %v\n", kubeContext)
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		clientcmd.NewDefaultClientConfigLoadingRules(),
		&clientcmd.ConfigOverrides{CurrentContext: kubeContext}).ClientConfig()
	if err != nil {
		panic(err.Error())
	}

	// Creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func patchServiceWithLabel(clientSet *kubernetes.Clientset, serviceName string, connected bool) error {
	ctx := context.TODO()
	if !connected {
		log.Debugf("removing true from %v\n", serviceName)
		payload := []patchStringValue{{
			Op:    "remove",
			Path:  "/spec/selector/kubetunnel",
			Value: "true",
		}}
		payloadBytes, _ := json.Marshal(payload)
		_, err := clientSet.
			CoreV1().
			Services("default").
			Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
		return err
	} else {
		log.Debugf("adding true to %v\n", serviceName)
		payload := []patchStringValue{{
			Op:    "add",
			Path:  "/spec/selector/kubetunnel",
			Value: "true",
		}}
		payloadBytes, _ := json.Marshal(payload)
		_, err := clientSet.
			CoreV1().
			Services("default").
			Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
		return err
	}
}

type LogOutputSplitter struct{}

func (splitter *LogOutputSplitter) Write(p []byte) (n int, err error) {
	if bytes.Contains(p, []byte("level=error")) || bytes.Contains(p, []byte("level=warn")) {
		return os.Stderr.Write(p)
	}
	return os.Stdout.Write(p)
}
