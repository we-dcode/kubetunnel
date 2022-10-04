package tomlutil

import (
	"github.com/BurntSushi/toml"
	"strings"
)

func Marshal(v interface{}) (string, error)  {

	sb := strings.Builder{}

	err := toml.NewEncoder(&sb).Encode(&v)
	if err != nil {
		return "", err
	}

	return sb.String(), nil
}