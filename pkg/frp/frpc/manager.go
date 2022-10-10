package frpc

type Manager struct {
	Host string
	Port string
}

func NewManager(host, port string) *Manager {

	return &Manager{
		Host: host,
		Port: port,
	}
}

func (m Manager) RunFRPc() {

}
