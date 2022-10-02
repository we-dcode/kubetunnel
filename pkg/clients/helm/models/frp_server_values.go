package models

import "fmt"

type FRPServerValues struct {
	Ports             Ports              `json:"env_ports"`
	ServiceName       string             `json:"env_service_name,omitempty"`
	PodSelectorLabels []PodSelectorLabel `json:"pod_selector_labels,omitempty"`
}

func (v *FRPServerValues) KubeTunnelServiceName()  string{

	return fmt.Sprintf("kubetunnel-%s", v.ServiceName)
}