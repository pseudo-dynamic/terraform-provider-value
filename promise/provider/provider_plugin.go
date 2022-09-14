package provider

import (
	"os"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// Provider
func Provider() func() tfprotov6.ProviderServer {
	var logLevel string
	logLevel, ok := os.LookupEnv("TF_LOG")

	if !ok {
		logLevel = "info"
	}

	return func() tfprotov6.ProviderServer {
		return &(RawProviderServer{logger: hclog.New(&hclog.LoggerOptions{
			Level:  hclog.LevelFromString(logLevel),
			Output: os.Stderr,
		})})
	}
}
