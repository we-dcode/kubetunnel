package models

type FrpServerValues struct {
	Namespace         string             `yaml:"namespace,omitempty"`
	Ports             Ports              `yaml:"env_ports"`
	ServiceName       string             `yaml:"env_service_name,omitempty"`
	PodSelectorLabels []PodSelectorLabel `yaml:"pod_selector_labels,omitempty"`
}
