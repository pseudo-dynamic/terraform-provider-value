package main

import (
	"context"
	"log"
	"os"

	promise "github.com/pseudo-dynamic/terraform-provider-value/promise/provider"
	stash "github.com/pseudo-dynamic/terraform-provider-value/stash/provider"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	"github.com/hashicorp/terraform-plugin-mux/tf5muxserver"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	stashProvider := stash.Provider()
	lazyProvider := promise.Provider()

	ctx := context.Background()
	muxer, err := tf5muxserver.NewMuxServer(ctx, stashProvider, lazyProvider)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	tf5server.Serve("registry.terraform.io/pseudo-dynamic/value", muxer.ProviderServer)
}
