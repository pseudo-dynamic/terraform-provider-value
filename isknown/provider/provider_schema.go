package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// GetProviderSchema function
func (s *UserProviderServer) GetProviderSchema(ctx context.Context, req *tfprotov6.GetProviderSchemaRequest) (*tfprotov6.GetProviderSchemaResponse, error) {
	return &tfprotov6.GetProviderSchemaResponse{
		Provider:        getProviderSchema(),
		ResourceSchemas: getDocumentedProviderResourceSchema(s.resourceSchemaParams),
		ProviderMeta:    common.GetProviderMetaSchema(),
	}, nil
}

func getProviderSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Version: 0,
		Block:   &tfprotov6.SchemaBlock{},
	}
}
