package main

import (
	"context"
	"log"
	"os"

	internal "github.com/pseudo-dynamic/terraform-provider-value/internal/provider"
	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
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

	isKnownProvider := isknown.ProviderConstructor(isknown.ProviderParameters{
		CheckFullyKnown: false,
	}, common.ProviderResourceSchemaParameters{
		ResourceName: "value_is_known",
		ResourceDescription: "Allows you to have a access to `result` during plan phase that " +
			"states whether `value` marked as \"(known after apply)\" or not.",
		ValueDescription: "The `value` (not nested attributes) is test against \"(known after apply)\"",
		ResultDescription: "States whether `value` is marked as \"(known after apply)\" or not. If `value` is an aggregate " +
			"type, only the top level of the aggregate type is checked; elements and attributes " +
			"are not checked.",
	})

	isFullyKnownProvider := isknown.ProviderConstructor(isknown.ProviderParameters{
		CheckFullyKnown: true,
	}, common.ProviderResourceSchemaParameters{
		ResourceName: "value_is_fully_known",
		ResourceDescription: "Allows you to have a access to `result` during plan phase that " +
			"states whether `value` or any nested attribute is marked as \"(known after apply)\" or not.",
		ValueDescription: "The `value` and if existing, nested attributes, are tested against \"(known after apply)\"",
		ResultDescription: "States whether `value` or any nested attribute is marked as \"(known after apply)\" or not. If `value` is an aggregate " +
			"type, not only the top level of the aggregate type is checked; elements and attributes " +
			"are checked too.",
	})

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
