package models

import (
	"fmt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)

type KubeTunnelResource struct {
	metav1.TypeMeta
	Metadata metav1.ObjectMeta      `json:"metadata" yaml:"metadata"`
	Spec     KubeTunnelResourceSpec `json:"spec" yaml:"spec"`
}

type KubeTunnelResourceSpec struct {
	Ports       Ports             `json:"env_ports" yaml:"env_ports"`
	ServiceName string            `json:"env_service_name,omitempty" yaml:"env_service_name,omitempty"`
	PodLabels   map[string]string `json:"pod_labels,omitempty" yaml:"pod_labels,omitempty"`
	Proxies     map[string]int    `json:"proxy,omitempty" yaml:"proxy,omitempty"`
}

func (v *KubeTunnelResourceSpec) KubeTunnelServiceName() string {

	return fmt.Sprintf("kubetunnel-%s", v.ServiceName)
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

func (p *Ports) UnmarshalJSON(b []byte) error {

	s := string(b)

	values := strings.Split(s, ",")

	for _, value := range values {
		p.Values = append(p.Values, strings.Trim(value, "\""))
	}

	return nil
}

func (p Ports) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf("\"%s\"", p.String())), nil
}
