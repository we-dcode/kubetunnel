package models

import (
	"fmt"
	"strings"
)

type Ports struct {
	Values []string
}

func (p Ports) String() string {

	return strings.Join(p.Values, ",")
}

func (p Ports) MarshalYAML() (interface{}, error) {

	return fmt.Sprintf("\"%s\"", p.String()), nil
}
