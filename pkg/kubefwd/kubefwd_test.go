package kubefwd_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/clients/kube"
	"github.com/we-dcode/kube-tunnel/pkg/kubefwd"
	"testing"
)

// TODO: execute this test using sudo -E go test -run 'KubeFwd'. the test doesn't work from IDE.
func TestKubeFwd(t *testing.T) {

	kubeClient := kube.MustNew("")

	err := kubefwd.Execute(kubeClient)

	assert.NoError(t, err)
}