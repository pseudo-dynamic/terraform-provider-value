package provider

import (
	"context"
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5/tf5server"
)

var providerName = "terraform-provider-value"

// Serve is the default entrypoint for the provider.
func Serve(ctx context.Context, logger hclog.Logger) error {
	return tf5server.Serve(providerName, func() tfprotov5.ProviderServer { return &(RawProviderServer{logger: logger}) })
}

// Provider
func Provider() func() tfprotov5.ProviderServer {
	var logLevel string
	logLevel, ok := os.LookupEnv("TF_LOG")

	if !ok {
		logLevel = "info"
	}

	return func() tfprotov5.ProviderServer {
		return &(RawProviderServer{logger: hclog.New(&hclog.LoggerOptions{
			Level:  hclog.LevelFromString(logLevel),
			Output: os.Stderr,
		})})
	}
}
