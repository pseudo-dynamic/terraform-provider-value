package goproviderconfig

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"
)

// GetProviderConfigType returns the tftypes.Type of a resource of type 'name'
func GetProviderConfigType() tftypes.Type {
	sch := GetProviderConfigSchema()
	return schema.GetObjectTypeFromSchema(sch)
}

func GetProviderConfigSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Version: 0,
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				GetGuidSeedAdditionSchemaAttribute(GetGuidSeedAdditionAttributeDescription()),
			},
		},
	}
}

func TryExtractProviderConfigGuidSeedAddition(providerConfig *tfprotov6.DynamicValue) (string, []*tfprotov6.Diagnostic, bool) {
	providerConfigType := GetProviderConfigType()
	return TryUnmarshalDynamicValueThenExtractGuidSeedAddition(providerConfig, providerConfigType)
}
