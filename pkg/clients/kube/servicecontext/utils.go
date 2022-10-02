package servicecontext

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/we-dcode/kube-tunnel/pkg/clients/helm/models"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	frpmodels "github.com/we-dcode/kube-tunnel/pkg/frp/models"
	v1 "k8s.io/api/core/v1"
	"strconv"
)

func ToFRPServerValues(ctx *ServiceContext) *models.FRPServerValues {

	var ports []string

	//var labelSelectors []models.PodSelectorLabel

	linq.From(ctx.Ports).Select(func(kubePort interface{}) interface{} {
		return strconv.Itoa(int(kubePort.(v1.ServicePort).Port))
	}).ToSlice(&ports)

	//linq.From(ctx.LabelSelector).Select(func(labelSelector interface{}) interface{} {
	//
	//	kv := labelSelector.(linq.KeyValue)
	//
	//	return models.PodSelectorLabel{
	//		Key:   kv.Key.(string),
	//		Value: kv.Value.(string),
	//	}
	//
	//}).ToSlice(&labelSelectors)

	return &models.FRPServerValues{
		Ports:             models.Ports{Values: ports},
		ServiceName:       ctx.ServiceName,
		PodSelectorLabels: ctx.LabelSelector,
	}
}

func ToFRPClientPairs(localIP string, remotePortByLocal map[string]string, ctx *ServiceContext) []frpc.ServicePair {

	var servicePairs []frpc.ServicePair

	localPortByRemote := make(map[string]string)

	for k, v := range remotePortByLocal {
		localPortByRemote[v] = k
	}

	linq.From(ctx.Ports).Select(func(kubePort interface{}) interface{} {

		port := strconv.Itoa(int(kubePort.(v1.ServicePort).Port))

		localPort := localPortByRemote[port]

		if localPort == "" {
			return nil // port not found in map
		}

		return frpc.ServicePair{
			Name: fmt.Sprintf("%s-%s", ctx.ServiceName, port),
			Service: frpmodels.Service{
				Type:       "tcp",
				RemotePort: port,
				LocalIP:    localIP,
				LocalPort:  localPort, //TODO: check what happen if port is not found in map
			},
		}

	}).Where(func(servicePair interface{}) bool {
		return servicePair != nil
	}).ToSlice(&servicePairs)

	return servicePairs
}
