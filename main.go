package main

import (
	"context"
	"log"
	"os"

	stash "github.com/teneko/terraform-provider-value/stash/provider"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
	tfmux "github.com/hashicorp/terraform-plugin-mux"
)

// Generate the Terraform provider documentation using `tfplugindocs`:
//go:generate go run github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs

func main() {
	valueProvider := stash.Provider()

	ctx := context.Background()
	factory, err := tfmux.NewSchemaServerFactory(ctx, valueProvider)

	if err != nil {
		log.Println(err.Error())
		os.Exit(1)
	}

	tf5server.Serve("registry.terraform.io/teneko/value", func() tfprotov5.ProviderServer {
		return factory.Server()
	})
}
