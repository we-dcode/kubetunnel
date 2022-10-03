package models

import "fmt"

type FRPServerValues struct {
	Ports             Ports             `yaml:"env_ports"`
	ServiceName       string            `yaml:"env_service_name,omitempty"`
	PodSelectorLabels map[string]string `yaml:"pod_selector_labels,omitempty"`
}

func (v *FRPServerValues) KubeTunnelServiceName() string {

	return fmt.Sprintf("kubetunnel-%s", v.ServiceName)
}
