package common

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
)

type ProviderResourceSchemaParameters struct {
	ResourceName        string
	ResourceDescription string
	ValueDescription    string
	ResultDescription   string
}

// GetProviderResourceSchema contains the definitions of all supported resources
func GetProviderResourceSchema(schema ProviderResourceSchemaParameters) map[string]*tfprotov6.Schema {
	return map[string]*tfprotov6.Schema{
		schema.ResourceName: {
			Version: 1,
			Block: &tfprotov6.SchemaBlock{
				Description: schema.ResourceDescription + "\n" + goproviderconfig.GetProviderMetaGuidSeedAdditionAttributeDescription(),
				BlockTypes:  []*tfprotov6.SchemaNestedBlock{},
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name:        "value",
						Type:        tftypes.DynamicPseudoType,
						Required:    true,
						Optional:    false,
						Computed:    false,
						Description: schema.ValueDescription,
					},
					{
						Name:        "guid_seed",
						Type:        tftypes.String,
						Required:    true,
						Optional:    false,
						Computed:    false,
						Description: goproviderconfig.GetGuidSeedAttributeDescription(schema.ResourceName),
					},
					{
						Name:        "proposed_unknown",
						Type:        tftypes.DynamicPseudoType,
						Required:    true,
						Optional:    false,
						Computed:    false,
						Description: goproviderconfig.GetProposedUnknownAttributeDescription(),
					},
					{
						Name:        "result",
						Type:        tftypes.Bool,
						Required:    false,
						Optional:    true,
						Computed:    true,
						Description: schema.ResultDescription,
					},
				},
			},
		},
	}
}
