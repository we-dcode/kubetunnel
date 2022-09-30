package hostsutils_test

import (
	"github.com/txn2/txeh"
	"github.com/we-dcode/kube-tunnel/pkg/utils/hostsutils"
	"gotest.tools/v3/assert"
	"testing"
)

func TestReplaceHostAddress(t *testing.T) {

	hostFile, err := txeh.NewHostsDefault()

	assert.NilError(t, err)

	hostsutils.ReplaceAddressForHost(hostFile, "nginx", "kubetunnel-nginx")
}
