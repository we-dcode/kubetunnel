/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/operator-framework/helm-operator-plugins/pkg/annotation"
	"github.com/operator-framework/helm-operator-plugins/pkg/reconciler"
	"github.com/operator-framework/helm-operator-plugins/pkg/watches"
	ctrlruntime "k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	//+kubebuilder:scaffold:imports
)

var (
	scheme                         = ctrlruntime.NewScheme()
	setupLog                       = ctrl.Log.WithName("setup")
	defaultMaxConcurrentReconciles = runtime.NumCPU()
	defaultReconcilePeriod         = time.Minute
)

type ServiceReq struct {
	Name        string
	Namespace   string
	Client      string
	IsConnected string
}

func init() {
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))

	//+kubebuilder:scaffold:scheme
}

func main() {
	var (
		metricsAddr          string
		leaderElectionID     string
		watchesPath          string
		probeAddr            string
		enableLeaderElection bool
	)

	setupLog.Info("Starting Gin server..")
	go startGin()

	// #### Manager pod configuration

	flag.StringVar(&metricsAddr, "metrics-bind-address", ":8080", "The address the metric endpoint binds to.")
	flag.StringVar(&probeAddr, "health-probe-bind-address", ":8081", "The address the probe endpoint binds to.")
	flag.StringVar(&watchesPath, "watches-file", "watches.yaml", "path to watches file")
	flag.StringVar(&leaderElectionID, "leader-election-id", "5dbd2493.dcode.tech", "provide leader election")
	flag.BoolVar(&enableLeaderElection, "leader-elect", false,
		"Enable leader election for controller manager. "+
			"Enabling this will ensure there is only one active controller manager.")
	opts := zap.Options{
		Development: true,
	}
	opts.BindFlags(flag.CommandLine)
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&opts)))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     metricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: probeAddr,
		LeaderElection:         enableLeaderElection,
		LeaderElectionID:       leaderElectionID,
	})
	if err != nil {
		setupLog.Error(err, "unable to start manager")
		os.Exit(1)
	}

	//+kubebuilder:scaffold:builder

	if err := mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up health check")
		os.Exit(1)
	}
	if err := mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		setupLog.Error(err, "unable to set up ready check")
		os.Exit(1)
	}

	ws, err := watches.Load(watchesPath)
	if err != nil {
		setupLog.Error(err, "Failed to create new manager factories")
		os.Exit(1)
	}

	for _, w := range ws {
		// Register controller with the factory
		reconcilePeriod := defaultReconcilePeriod
		if w.ReconcilePeriod != nil {
			reconcilePeriod = w.ReconcilePeriod.Duration
		}

		maxConcurrentReconciles := defaultMaxConcurrentReconciles
		if w.MaxConcurrentReconciles != nil {
			maxConcurrentReconciles = *w.MaxConcurrentReconciles
		}

		// Setup manager with Helm API

		w.OverrideValues = map[string]string{}
		w.OverrideValues["operator.namespace"] = getEnvVar("POD_NAMESPACE")
		w.OverrideValues["operator.service.port"] = getEnvVar("SERVICE_PORT")
		w.OverrideValues["operator.service.name"] = getEnvVar("SERVICE_NAME")
		fmt.Printf("Values: %+v", w.OverrideValues)

		r, err := reconciler.New(
			reconciler.WithChart(*w.Chart),
			reconciler.WithGroupVersionKind(w.GroupVersionKind),
			reconciler.WithOverrideValues(w.OverrideValues),
			reconciler.SkipDependentWatches(w.WatchDependentResources != nil && !*w.WatchDependentResources),
			reconciler.WithMaxConcurrentReconciles(maxConcurrentReconciles),
			reconciler.WithReconcilePeriod(reconcilePeriod),
			reconciler.WithInstallAnnotations(annotation.DefaultInstallAnnotations...),
			reconciler.WithUpgradeAnnotations(annotation.DefaultUpgradeAnnotations...),
			reconciler.WithUninstallAnnotations(annotation.DefaultUninstallAnnotations...),
		)
		if err != nil {
			setupLog.Error(err, "unable to create helm reconciler", "controller", "Helm")
			os.Exit(1)
		}

		// The SetupWithManager method is called when the operator starts.
		// It serves to tell the operator framework what types our PodReconciler
		// needs to watch. To use the same Pod type used by Kubernetes internally,
		// we need to import some of its code. All of the Kubernetes source code is
		// open source, so you can import any part you like in your own Go code. You can
		// find a complete list of available packages in the Kubernetes source code or here
		// on pkg.go.dev. To use pods, we need the k8s.io/api/core/v1 package.
		if err := r.SetupWithManager(mgr); err != nil {

			setupLog.Error(err, "unable to create controller", "controller", "Helm")
			os.Exit(1)
		}
		setupLog.Info("configured watch", "gvk", w.GroupVersionKind, "chartPath", w.ChartPath, "maxConcurrentReconciles", maxConcurrentReconciles, "reconcilePeriod", reconcilePeriod)
	}

	setupLog.Info("starting manager")
	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		setupLog.Error(err, "problem running manager")
		os.Exit(1)
	}
}

func startGin() {
	setupLog.Info("Inside Gin function..")
	router := gin.Default()
	router.POST("/service", serviceHandler)

	router.Run("0.0.0.0:8083")
	setupLog.Info("Started Gin Server on port 8083..")
}

func serviceHandler(c *gin.Context) {
	var req ServiceReq
	c.BindJSON(&req)
	c.JSON(200, req)

	kube := connectToKubernetes(req.Namespace) // TODO: change this one
	isConnected, err := strconv.ParseBool(req.IsConnected)
	if err != nil {
		fmt.Errorf(err.Error())
	}
	error := patchServiceWithLabel(kube, req.Name, isConnected)
	if error != nil {
		fmt.Errorf("error %v", err)
	} else {
		fmt.Print("No error!")
	}

}

func connectToKubernetes(namespace string) *kube.Kube {

	kubeClient := kube.MustNew("~/.kube/config", namespace)

	err := kubeClient.ConnectivityCheck()
	if err != nil {
		fmt.Printf(err.Error())
	}

	return kubeClient
}

func patchServiceWithLabel(k *kube.Kube, serviceName string, connected bool) error {

	svcContext, err := k.GetServiceContext(serviceName)
	if err != nil {
		log.Errorf("fail to get service: '%s' context, error: %s", serviceName, err.Error())
	}

	clientSet := k.InnerKubeClient
	log.Debugf(k.Namespace)
	ctx := context.TODO()

	slugPrefix := fmt.Sprintf("%s-", constants.KubetunnelSlug)

	// TODO: Check if service needs to be updated...

	if !connected {

		log.Debugf("removing true from %v\n", serviceName)
		var payload []kube.PatchOperation

		for key, valueWithSlug := range svcContext.LabelSelector {

			if strings.EqualFold(key, constants.KubetunnelSlug) {
				continue
			}

			slugAlreadyRemoved := strings.Contains(valueWithSlug, constants.KubetunnelSlug) == false
			if slugAlreadyRemoved {
				continue
			}

			valueWithoutSlug := strings.Replace(valueWithSlug, slugPrefix, "", 1)

			payload = append(payload, kube.PatchOperation{

				Op:    "replace",
				Path:  fmt.Sprintf("/spec/selector/%s", key),
				Value: valueWithoutSlug,
			})
		}

		if len(payload) > 0 {
			payloadBytes, _ := json.Marshal(payload)
			_, err := clientSet.
				CoreV1().
				Services(k.Namespace).
				Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
			return err
		}

	} else {
		log.Debugf("adding true to %v\n", serviceName)

		var payload []kube.PatchOperation

		for key, valueWithoutSlug := range svcContext.LabelSelector {

			if strings.EqualFold(key, constants.KubetunnelSlug) {
				continue
			}

			alreadyContainSlug := strings.Contains(valueWithoutSlug, constants.KubetunnelSlug)

			if alreadyContainSlug {
				continue
			}

			valueWithSlug := fmt.Sprintf("%s%s", slugPrefix, valueWithoutSlug)

			payload = append(payload, kube.PatchOperation{

				Op:    "replace",
				Path:  fmt.Sprintf("/spec/selector/%s", key),
				Value: valueWithSlug,
			})
		}
		if len(payload) > 0 {
			payloadBytes, _ := json.Marshal(payload)
			_, err := clientSet.
				CoreV1().
				Services(k.Namespace).
				Patch(ctx, serviceName, types.JSONPatchType, payloadBytes, metav1.PatchOptions{})
			return err
		}
	}

	return nil

}

func getEnvVar(variable string) string {
	envVar, ok := os.LookupEnv(variable)
	if !ok {
		fmt.Errorf("%v is not a present env variable", variable)
	}
	return envVar
}
