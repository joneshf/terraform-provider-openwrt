package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/joneshf/terraform-provider-openwrt/openwrt"
)

// Provider documentation generation.
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs generate --provider-name openwrt

func main() {
	ctx := context.Background()
	providerNew := openwrt.New
	options := providerserver.ServeOpts{
		Address: "registry.terraform.io/joneshf/terraform-provider-openwrt",
	}
	providerserver.Serve(
		ctx,
		providerNew,
		options,
	)
}
