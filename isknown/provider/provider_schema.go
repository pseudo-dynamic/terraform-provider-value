package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
)

// GetProviderSchema function
func (s *UserProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return &tfprotov6.GetProviderSchemaResponse{
		Provider:        goproviderconfig.GetProviderConfigSchema(),
		ResourceSchemas: getDocumentedProviderResourceSchema(s.resourceSchemaParams),
		ProviderMeta:    goproviderconfig.GetProviderMetaSchema(),
	}, nil
}
