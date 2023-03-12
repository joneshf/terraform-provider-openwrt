package main

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name openwrt

const (
	// version is set by the linker flag `-X main.version=<some-version>`.
	version = "unknown"
)

func main() {
	ctx := context.Background()
	providerNew := func() provider.Provider {
		return openwrt.New(version, os.LookupEnv)
	}
	options := providerserver.ServeOpts{
		Address: "registry.terraform.io/joneshf/openwrt",
	}
	providerserver.Serve(
		ctx,
		providerNew,
		options,
	)
}
