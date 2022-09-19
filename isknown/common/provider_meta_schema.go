package common

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"
)

// GetProviderMetaType returns the tftypes.Type of a resource of type 'name'
func GetProviderMetaType() tftypes.Type {
	sch := GetProviderMetaSchema()
	return schema.GetObjectTypeFromSchema(sch)
}

func GetProviderMetaSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Version: 0,
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				goproviderconfig.GetGuidSeedAdditionSchemaAttribute(getProviderMetaGuidSeedAdditionAttributeDescription()),
			},
		},
	}
}

func getProviderMetaGuidSeedAdditionAttributeDescription() string {
	return "## Provider Metadata\n" +
		"Each module can use provider_meta. Please keep in mind that these settings only count " +
		"for resources of this module! (see [https://www.terraform.io/internals/provider-meta](https://www.terraform.io/internals/provider-meta)):\n" +
		"```terraform\n" + `// Terraform provider_meta example
terraform {
	// "value" is the provider name
	provider_meta "value" {
		// {workdir} -> The only available placeholder currently (see below for more information)
		guid_seed_addition = "{workdir}#for-example" // Results into "/path/to/workdir#for-example"
	}
}` + "\n```\n" +
		"### Optional\n" +
		"- `guid_seed_addition` (String) " + goproviderconfig.GetGuidSeedAdditionAttributeDescription()
}

func TryExtractProviderMetaGuidSeedAddition(providerMeta *tfprotov6.DynamicValue) (string, []*tfprotov6.Diagnostic, bool) {
	providerMetaType := GetProviderMetaType()
	return goproviderconfig.TryExtractGetSeedAddition(providerMeta, providerMetaType)
}
