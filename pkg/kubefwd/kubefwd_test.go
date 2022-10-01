package kubefwd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/kubefwd"
	"testing"
)

func TestKubeFwd(t *testing.T) {

	kubeClient := kube.MustNew("")

	err := kubefwd.Execute(kubeClient)

	assert.NoError(t, err)
}