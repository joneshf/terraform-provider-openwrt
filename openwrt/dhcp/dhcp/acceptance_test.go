//go:build acceptance.test

package dhcp_test

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/joneshf/terraform-provider-openwrt/internal/acceptancetest"
	"github.com/ory/dockertest/v3"
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
