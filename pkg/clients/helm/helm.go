package helm

import (
	"context"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm/models"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"gopkg.in/yaml.v3"
	"helm.sh/helm/v3/pkg/repo"
)

type Helm struct {
	namespace string
	helmClient helmclient.Client
}

type EMPTY struct {

}

func MustNew(kube *kube.Kube) *Helm {

	client, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace:        kube.Namespace,
		},
		RestConfig: kube.Config,
	})

	if err != nil {
		log.Panicf("fail creating helm client, more info: '%s'", err.Error())
	}

	err = client.AddOrUpdateChartRepo(repo.Entry{
		Name: "we-decode",
		URL:  constants.DcodeChartRepo,
	})

	if err != nil {
		log.Panicf("fail adding Dcode's chart repo, more info: '%s'", err.Error())
	}

	return &Helm{
		kube.Namespace,
		client,
	}
}


func (c *Helm) InstallOrUpgradeFrpServer(chartVersion string, values *models.FRPServerValues) error {

	releaseName := fmt.Sprintf("kubetunnel-%s", values.ServiceName)

	return install(c, constants.KubeTunnelChartName, chartVersion, releaseName, values)
}


func (c *Helm) InstallOrUpgradeGC(chartVersion string) error{

	releaseName := "kubetunnel-gc"

	return install(c, constants.KubeTunnelChartName, chartVersion, releaseName, EMPTY{})
}

func install(c *Helm, chartName string, chartVersion string, releaseName string, values interface{}) error {
	valuesYaml, err := yaml.Marshal(values)
	if err != nil {
		return fmt.Errorf("err: fail to parse values.yaml more info: '%s'", err.Error())
	}

	chartSpec := helmclient.ChartSpec{
		ReleaseName: releaseName,
		Recreate:    true,
		ChartName:   chartName,
		//Version:     chartVersion,
		Namespace:   c.namespace,
		UpgradeCRDs: false,
		Wait:        true,
		ValuesYaml:  string(valuesYaml),
	}

	if _, err := c.helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil); err != nil {
		return fmt.Errorf("err: fail installing %s chart. more info: %s", chartName, err.Error())
	}

	return nil
}
