package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetProviderSchema function
func (s *RawProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov5.GetProviderSchemaRequest) (*tfprotov5.GetProviderSchemaResponse, error) {
	cfgSchema := GetProviderConfigSchema()
	resSchema := GetProviderResourceSchema()

	log.Println("--------------------------GetProviderSchema Called------------------------------")

	return &tfprotov5.GetProviderSchemaResponse{
		Provider:        cfgSchema,
		ResourceSchemas: resSchema,
	}, nil
}
