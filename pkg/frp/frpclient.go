package frp

import(
	"github.com/fatedier/frp/client"
	"github.com/fatedier/frp/pkg/config"
)

type FrpClient struct {

}

func New() *FrpClient {

	tcpProxyConf := &config.TCPProxyConf{
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

	tcpConf := map[string]config.ProxyConf{
		tcpProxyConf.ProxyName: tcpProxyConf,
	}
	var cfg config.ClientCommonConf
// https://github.com/fatedier/frp/blob/6ecc97c8571df002dd7cf42522e3f2ce9de9a14d/cmd/frpc/sub/root.go#L200
	svr, errRet := client.NewService(cfg, tcpConf, nil, "")
	_,_ = svr, errRet
	//if errRet != nil {
	//	err := errRet
	//	return
	//}
	//
	//kcpDoneCh := make(chan struct{})
	//// Capture the exit signal if we use kcp.
	//if cfg.Protocol == "kcp" {
	//	go handleSignal(svr, kcpDoneCh)
	//}
	//
	//err = svr.Run()
	//if err == nil && cfg.Protocol == "kcp" {
	//	<-kcpDoneCh
	//}
	//return
	//
	//return nil
}

func (f *FrpClient) Execute()  {

}
