package provider

import (
	"context"
	"log"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetProviderSchema function
func (s *RawProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	cfgSchema := goproviderconfig.GetProviderConfigSchema()
	resSchema := GetProviderResourceSchema()

	log.Println("--------------------------GetProviderSchema Called------------------------------")

	return &tfprotov6.GetProviderSchemaResponse{
		Provider:        cfgSchema,
		ResourceSchemas: resSchema,
	}, nil
}
