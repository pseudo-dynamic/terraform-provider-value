package main

import (
	"context"
	"log"
	"os"

	internal "github.com/pseudo-dynamic/terraform-provider-value/internal/fwkprovider"
	isfullyknown "github.com/pseudo-dynamic/terraform-provider-value/isfullyknown/provider"
	isknown "github.com/pseudo-dynamic/terraform-provider-value/isknown/provider"
	promise "github.com/pseudo-dynamic/terraform-provider-value/promise/provider"
	stash "github.com/pseudo-dynamic/terraform-provider-value/stash/provider"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6/tf6server"
	"github.com/hashicorp/terraform-plugin-mux/tf6muxserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	stashProvider := stash.Provider()
	promiseProvider := promise.Provider()
	isKnownProvider := isknown.Provider()
	isFullyKnownProvider := isfullyknown.Provider()
	internalProvider := providerserver.NewProtocol6(internal.NewProvider())
	ctx := context.Background()

	muxer, err := tf6muxserver.NewMuxServer(
		ctx,
		stashProvider,
		promiseProvider,
		isKnownProvider,
		isFullyKnownProvider,
		internalProvider,
	)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	tf6server.Serve("registry.terraform.io/pseudo-dynamic/value", muxer.ProviderServer)
}
