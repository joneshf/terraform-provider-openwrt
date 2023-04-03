//go:build acceptance.test

package wifidevice_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/joneshf/terraform-provider-openwrt/lucirpc"
	"github.com/ory/dockertest/v3"
	"golang.org/x/crypto/ssh"
	"gotest.tools/v3/assert"
)

var (
	dockerPool *dockertest.Pool
)

func TestMain(m *testing.M) {
	var (
		code     int
		err      error
		tearDown func()
	)
	ctx := context.Background()
	tearDown, dockerPool, err = acceptancetest.Setup(ctx)
	defer func() {
		tearDown()
		os.Exit(code)
	}()
	if err != nil {
		fmt.Printf("Problem setting up tests: %s", err)
		code = 1
		return
	}

	log.Printf("Running tests")
	code = m.Run()
}

// runOpenWrtServerWithWireless starts an OpenWrt server,
// and sets up the wireless config.
// Without setting up the config,
// the tests in this package will fail.
func runOpenWrtServerWithWireless(
	ctx context.Context,
	dockerPool dockertest.Pool,
	t *testing.T,
) (*lucirpc.Client, string) {
	openWrtServer := acceptancetest.RunOpenWrtServer(
		ctx,
		dockerPool,
		t,
	)
	sshURL := fmt.Sprintf("%s:%d", openWrtServer.Hostname, openWrtServer.SSHPort)
	sshConfig := &ssh.ClientConfig{
		Auth: []ssh.AuthMethod{
			ssh.Password(openWrtServer.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		User:            openWrtServer.Username,
	}
	sshClient, err := ssh.Dial("tcp", sshURL, sshConfig)
	assert.NilError(t, err)
	t.Cleanup(func() {
		sshClient.Close()
	})
	session, err := sshClient.NewSession()
	assert.NilError(t, err)
	t.Cleanup(func() {
		session.Close()
	})
	err = session.Run("touch /etc/config/wireless")
	assert.NilError(t, err)
	client := openWrtServer.LuCIRPCClient(
		ctx,
		t,
	)
	providerBlock := openWrtServer.ProviderBlock()
	return client, providerBlock
}
