package models

import "fmt"

type PodSelectorLabel struct {

	Key string
	Value string
}

func (p PodSelectorLabel) String() string  {

	return fmt.Sprintf("%s:%s", p.Key, p.Value)
}
