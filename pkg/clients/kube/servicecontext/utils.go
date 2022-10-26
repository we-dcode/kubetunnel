package servicecontext

import (
	"fmt"
	"github.com/ahmetb/go-linq/v3"
	"github.com/we-dcode/kube-tunnel/pkg/constants"
	"github.com/we-dcode/kube-tunnel/pkg/frp/frpc"
	frpmodels "github.com/we-dcode/kube-tunnel/pkg/frp/models"
	kubeTunnelModels "github.com/we-dcode/kube-tunnel/pkg/models"
	v1 "k8s.io/api/core/v1"
	"strconv"
)

func ToKubeTunnelResourceSpec(ctx *ServiceContext, podLabels map[string]string) kubeTunnelModels.KubeTunnelResourceSpec {

	var ports []string

	linq.From(ctx.Ports).Select(func(kubePort interface{}) interface{} {
		return strconv.Itoa(int(kubePort.(v1.ServicePort).Port))
	}).ToSlice(&ports)

	//labelSelectors := make(map[string]string)
	//
	//for key, value := range ctx.LabelSelector {
	//
	//	labelSelectors[key] = fmt.Sprintf("%s-%s", constants.KubetunnelSlug, value)
	//}

	podLabels[constants.KubetunnelSlug] = ctx.ServiceName

	return kubeTunnelModels.KubeTunnelResourceSpec{
		Ports:       kubeTunnelModels.Ports{Values: ports},
		ServiceName: ctx.ServiceName,
		PodLabels:   podLabels,
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
