package helm

import (
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/models"
)

type Helm struct {

}

func NewHelmClient()  {

	//helmclient.NewClientFromKubeConf(config, helmclient.RestConfClientOptions{
	//	Options:    nil,
	//	RestConfig: nil,
	//})
	//
	//helmClient, err := helmclient.New(&helmclient.Options{
	//	Debug:            false,
	//	Linting:          false,
	//	DebugLog:         nil,
	//	RegistryConfig:   "",
	//	Output:           nil,
	//})
}


func (c *Helm) MustInstallFrpServer(chartUrl string, values *models.FrpServerValues)  {

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



	// Install a chart release.
	// Note that helmclient.Options.Namespace should ideally match the namespace in chartSpec.Namespace.
	if _, err := helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil); err != nil {
		log.Panicf("err: fail installing frp chart. more info: %s", err.Error())
	}
}


func (c *Helm) InstallGarbageCollector()  {

}
