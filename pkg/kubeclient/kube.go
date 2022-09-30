package kubeclient

import (
	"context"
	"fmt"
	"github.com/ghodss/yaml"
	helmclient "github.com/mittwald/go-helm-client"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/models"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

type KubeClient struct {

	kubeClient kubernetes.Interface
}


func MustNew() *KubeClient  {

	kubeConfig, err := clientcmd.DefaultClientConfig.ClientConfig()
	if  err != nil {
		log.Panicf("err: unable to read kubeconfig. \"%s\"", err.Error())
	}

	kubeClient, err := kubernetes.NewForConfig(kubeConfig)
	if  err != nil {
		log.Panicf("err: unable to create kube client. \"%s\"", err.Error())
	}

	return &KubeClient{
		kubeClient,
	}
}

func (c *KubeClient) MustInstallFrpServer(chartUrl string, values *models.FrpServerValues)  {

	valuesYaml, err := yaml.Marshal(values)
	if err != nil {
		log.Panicf("err: fail to parse values.yaml more info: %s", err.Error())
	}

	if len(chartUrl) == 0 {
		chartUrl = constants.FrpServerChart
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName: fmt.Sprintf("kubetunnel-frpserver-%s", values.ServiceName),
		ChartName:   chartUrl,
		Namespace:   values.Namespace,
		UpgradeCRDs: true,
		Wait:        true,
		ValuesYaml: string(valuesYaml),
	}

	helmClient, _ := helmclient.New(&helmclient.Options{
		Namespace:        values.Namespace,
		RepositoryConfig: "",
		RepositoryCache:  "",
		Debug:            false,
		Linting:          false,
		DebugLog:         nil,
		RegistryConfig:   "",
		Output:           nil,
	})

	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil); err != nil {
		log.Panicf("err: fail installing frp chart. more info: %s", err.Error())
	}
}


func (c *KubeClient) InstallGarbageCollector()  {

}
