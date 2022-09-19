package provider

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// UserProviderServer implements the ProviderServer interface as exported from ProtoBuf.
type UserProviderServer struct {
	// Since the provider is essentially a gRPC server, the execution flow is dictated by the order of the client (Terraform) request calls.
	// Thus it needs a way to persist state between the gRPC calls. These attributes store values that need to be persisted between gRPC calls,
	// such as instances of the Kubernetes clients, configuration options needed at runtime.
	logger hclog.Logger
	//providerEnabled bool
	hostTFVersion              string
	params                     ProviderParameters
	resourceSchemaParams       common.ProviderResourceSchemaParameters
	ProviderConfigSeedAddition string
}

type ProviderParameters struct {
	CheckFullyKnown bool
}

// ProviderConstructor
func ProviderConstructor(providerParams ProviderParameters, resourceSchemaParams common.ProviderResourceSchemaParameters) func() tfprotov6.ProviderServer {
	var logLevel string
	logLevel, ok := os.LookupEnv("TF_LOG")

	if !ok {
		logLevel = "info"
	}

	return func() tfprotov6.ProviderServer {
		logger := hclog.New(&hclog.LoggerOptions{
			Level:  hclog.LevelFromString(logLevel),
			Output: os.Stderr,
		})

		return &UserProviderServer{
			logger:               logger,
			params:               providerParams,
			resourceSchemaParams: resourceSchemaParams,
		}
	}
}

func Provider() func() tfprotov6.ProviderServer {
	return ProviderConstructor(ProviderParameters{
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
}
