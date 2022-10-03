package models

type FrpClientConfig map[string]interface{}

//type FrpClientConfig struct {
//	Common Common `toml:"common,omitempty"`
//	Services map[string]Service `toml:"services,omitempty"`
//}

type Common struct {
	ServerAddress string `toml:"server_address,omitempty"`
	ServerPort    string `toml:"server_port,omitempty"`
}

type Service struct {
	Type       string `toml:"type,omitempty"`
	RemotePort string `toml:"remote_port,omitempty"`
	LocalIP    string `toml:"local_ip,omitempty"`
	LocalPort  string `toml:"local_port,omitempty"`
}
