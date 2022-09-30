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

func (p Ports) MarshalJSON() ([]byte, error)  {

	return []byte(fmt.Sprintf("\"%s\"", p.String())), nil
}