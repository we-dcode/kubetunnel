package tcputil_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/we-dcode/kube-tunnel/pkg/utils/tcputil"
	"testing"
)

func TestFRPConnectivityWhenFRPCDestinationIsNotAvailable(t *testing.T) {

	isAvailable := tcputil.IsAvailable("localhost", "22285")

	assert.True(t, isAvailable)
}
