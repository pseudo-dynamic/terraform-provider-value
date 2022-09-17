package common

import (
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// GetProviderMetaType returns the tftypes.Type of a resource of type 'name'
func GetProviderMetaType() tftypes.Type {
	sch := GetProviderMetaSchema()
	return getObjectTypeFromSchema(sch)
}

// getObjectTypeFromSchema returns a tftypes.Type that can wholy represent the schema input
// TODO: Outsource this and reduce redudancy across terraform-plugin-go providers
func getObjectTypeFromSchema(schema *tfprotov6.Schema) tftypes.Type {
	bm := map[string]tftypes.Type{}

	for _, att := range schema.Block.Attributes {
		bm[att.Name] = att.Type
	}

	for _, b := range schema.Block.BlockTypes {
		a := map[string]tftypes.Type{}
		for _, att := range b.Block.Attributes {
			a[att.Name] = att.Type
		}
		bm[b.TypeName] = tftypes.List{
			ElementType: tftypes.Object{AttributeTypes: a},
		}

		// FIXME we can make this function recursive to handle
		// n levels of nested blocks
		for _, bb := range b.Block.BlockTypes {
			aa := map[string]tftypes.Type{}
			for _, att := range bb.Block.Attributes {
				aa[att.Name] = att.Type
			}
			a[bb.TypeName] = tftypes.List{
				ElementType: tftypes.Object{AttributeTypes: aa},
			}
		}
	}

	return tftypes.Object{AttributeTypes: bm}
}

var seedPrefixDescription = "## Provider Metadata\n" +
	"Each module can use provider_meta. Please keep in mind that these settings only count " +
	"for resources of this module! (see https://www.terraform.io/internals/provider-meta):\n" +
	"```terraform\n" + `// Terraform provider_meta example
terraform {
	// "value" is the provider name
	provider_meta "value" {
		// {workdir} -> The only available placeholder currently (see below for more information)
		seed_prefix = "{workdir}#for-example" // Results into "/path/to/workdir#for-example"
	}
}` + "\n```\n" +
	"### Optional\n" +
	"- `seed_prefix` (String) It gets appended to each seed of any `value_is_fully_known` (resource) or " +
	"`value_is_known` (resource) within the same module.\n" + `
	**Placeholders**:
	- "{workdir}" (Keyword) The actual workdir; equals to terraform's path.root. This placeholder is
	recommended because this value won't be dragged along the plan and apply phase in comparison to
	"abspath(path.root)" that you would add to resource seed where a change to path.root would be
	recognized just as usual from terraform.`

func GetProviderMetaSchema() *tfprotov6.Schema {
	return &tfprotov6.Schema{
		Version: 0,
		Block: &tfprotov6.SchemaBlock{
			Attributes: []*tfprotov6.SchemaAttribute{
				{
					Name:        "seed_prefix",
					Type:        tftypes.String,
					Required:    false,
					Optional:    true,
					Computed:    false,
					Description: seedPrefixDescription,
				},
			},
		},
	}
}

func CanGetSeedPrefix(providerMeta *tfprotov6.DynamicValue) (string, []*tfprotov6.Diagnostic, bool) {
	var seedPrefixValue tftypes.Value
	var seedPrefix string

	var isErroneous bool
	var isWorking bool

	providerMetaType := GetProviderMetaType()
	providerMetaValueDynamic := providerMeta
	var providerMetaValue tftypes.Value
	var providerMetaValueMap map[string]tftypes.Value
	var diags []*tfprotov6.Diagnostic
	if providerMetaValue, providerMetaValueMap, diags, isErroneous = UnmarshalState(providerMetaValueDynamic, providerMetaType); isErroneous {
		goto End
	}
	_ = providerMetaValue

	if seedPrefixValue, isWorking = providerMetaValueMap["seed_prefix"]; !isWorking {
		goto End
	}
	_ = seedPrefixValue

	if !seedPrefixValue.IsKnown() {
		goto End
	}

	if err := seedPrefixValue.As(&seedPrefix); err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Extraction failed",
			Detail:   "Seed prefix could not be extracted as string",
		})

		goto End
	}

	{
		workdir, _ := os.Getwd()
		seedPrefix = strings.ReplaceAll(seedPrefix, "{workdir}", workdir)
	}

End:
	return seedPrefix, diags, len(diags) == 0
}
