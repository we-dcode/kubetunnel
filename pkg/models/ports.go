package models

import "strings"

type Ports struct {
	Values []string
}

func (p Ports) String() string {

	return strings.Join(p.Values, ",")
}

