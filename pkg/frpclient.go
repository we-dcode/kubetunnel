package pkg

import(
	"github.com/fatedier/frp/pkg/config"
)

type FrpClient struct {

}

func New() *FrpClient{

	tcpProxyConf := config.TCPProxyConf{
		BaseProxyConf: config.BaseProxyConf{
			ProxyName:            "",
			ProxyType:            "",
			UseEncryption:        false,
			UseCompression:       false,
			Group:                "",
			GroupKey:             "",
			ProxyProtocolVersion: "",
			BandwidthLimit:       config.BandwidthQuantity{},
			Metas:                nil,
			LocalSvrConf:         config.LocalSvrConf{},
			HealthCheckConf:      config.HealthCheckConf{},
		},
		RemotePort:    0,
	}


	return nil
}
