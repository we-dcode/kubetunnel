package main

import (
	"fmt"
	"testing"
)

func TestPatchServiceWithLabel(t *testing.T) {

	kube := connectToKubernetes("operator-system") // TODO: change this one

	error := patchServiceWithLabel(kube, "nginx", true)
	if error != nil {
		fmt.Errorf("error %v", error)
	} else {
		fmt.Print("No error!")
	}

}
