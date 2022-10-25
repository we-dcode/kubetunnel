package helm

import (
	"context"
	"fmt"
	helmclient "github.com/mittwald/go-helm-client"
	log "github.com/sirupsen/logrus"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"helm.sh/helm/v3/pkg/repo"
	"time"
)

type Helm struct {
	namespace  string
	helmClient helmclient.Client
}

type EMPTY struct {
}

func MustNew(kube *kube.Kube) *Helm {

	client, err := helmclient.NewClientFromRestConf(&helmclient.RestConfClientOptions{
		Options: &helmclient.Options{
			Namespace: kube.Namespace,
			Debug:     true,
			DebugLog: func(format string, v ...interface{}) {
				log.Infof(format, v)
			},
		},
		RestConfig: kube.Config,
	})

	if err != nil {
		log.Panicf("fail creating helm client, more info: '%s'", err.Error())
	}

	err = client.AddOrUpdateChartRepo(repo.Entry{
		Name: constants.DcodeSlug,
		URL:  constants.DcodeChartRepo,
	})

	if err != nil {
		log.Panicf("fail adding Dcode's chart repo, more info: '%s'", err.Error())
	}

	err = client.UpdateChartRepos()
	if err != nil {
		log.Panicf("fail updating Dcode's helm repo, more info: %s", err.Error())
	}

	return &Helm{
		kube.Namespace,
		client,
	}
}

//func (c *Helm) InstallOrUpgradeFrpServer(chartVersion string, values *models.FRPServerValues) error {
//
//	releaseName := values.KubeTunnelServiceName()
//
//	valuesYaml, err := yaml.Marshal(values)
//	if err != nil {
//		return fmt.Errorf("err: fail to parse values.yaml more info: '%s'", err.Error())
//	}
//
//	return install(c, constants.KubeTunnelChartName, chartVersion, releaseName, valuesYaml)
//}

func (c *Helm) InstallKubeTunnelOperator(chartVersion string) error {

	releaseName := "kubetunnel-operator"

	return installWithNamespace(c, constants.KubetunnelOperatorChartName, chartVersion, releaseName, constants.KubetunnelSlug, []byte{})
}

func install(c *Helm, chartName string, chartVersion string, releaseName string, valuesYaml []byte) error {

	return installWithNamespace(c, chartName, chartVersion, releaseName, c.namespace, valuesYaml)
}

func installWithNamespace(c *Helm, chartName string, chartVersion string, releaseName string, namespace string, valuesYaml []byte) error {

	chartSpec := helmclient.ChartSpec{
		ReleaseName:     releaseName,
		Recreate:        false,
		ChartName:       chartName,
		Atomic:          true,
		Version:         chartVersion,
		Namespace:       namespace,
		CreateNamespace: true,
		UpgradeCRDs:     false,
		Wait:            true,
		Replace:         false,
		ValuesYaml:      string(valuesYaml),
		Timeout:         time.Second * 30,
	}

	if _, err := c.helmClient.InstallOrUpgradeChart(context.Background(), &chartSpec, nil); err != nil {
		return fmt.Errorf("err: fail installing %s chart. more info: %s", chartName, err.Error())
	}

	return nil
}
