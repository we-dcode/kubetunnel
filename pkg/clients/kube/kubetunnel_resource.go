package kube

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type KubeTunnelResource struct {
	metav1.TypeMeta
	ObjectMeta metav1.ObjectMeta
	Spec       KubeTunnelResourceSpec
}

type KubeTunnelResourceSpec struct {
	Ports             Ports             `yaml:"env_ports"`
	ServiceName       string            `yaml:"env_service_name,omitempty"`
	PodSelectorLabels map[string]string `yaml:"pod_selector_labels,omitempty"`
}

type Ports struct {
	Values []string
}

func (p Ports) String() string {

	return strings.Join(p.Values, ",")
}

func (p Ports) MarshalYAML() (interface{}, error) {

	return fmt.Sprintf("\"%s\"", p.String()), nil
}
