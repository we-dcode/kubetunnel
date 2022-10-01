package models

type FRPServerValues struct {
	Ports             Ports              `json:"env_ports"`
	ServiceName       string             `json:"env_service_name,omitempty"`
	PodSelectorLabels []PodSelectorLabel `json:"pod_selector_labels,omitempty"`
}
