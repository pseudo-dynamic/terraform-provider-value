package goproviderconfig

import (
	"os"
	"strings"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func TryUnmarshalDynamicValueThenExtractGuidSeedAddition(dynamicValue *tfprotov6.DynamicValue, dynamicValueType tftypes.Type) (string, []*tfprotov6.Diagnostic, bool) {
	seedAddition := ""

	var isErroneous bool

	var stateValue tftypes.Value
	var stateValueMap map[string]tftypes.Value
	var diags []*tfprotov6.Diagnostic

	if stateValue, stateValueMap, diags, isErroneous = schema.UnmarshalDynamicValue(dynamicValue, dynamicValueType); isErroneous {
		goto Return
	}
	_ = stateValue
	_ = stateValueMap

	seedAddition, diags, _ = TryUnmarshalValueThenExtractGuidSeedAddition(&stateValue)

Return:
	return seedAddition, diags, len(diags) == 0
}

func TryUnmarshalValueThenExtractGuidSeedAddition(value *tftypes.Value) (string, []*tfprotov6.Diagnostic, bool) {
	var seedAdditionValue tftypes.Value
	seedAddition := ""

	var isErroneous bool
	var isSuccesful bool

	var valueMap map[string]tftypes.Value
	var diags []*tfprotov6.Diagnostic

	if valueMap, diags, isErroneous = schema.UnmarshalValue(value); isErroneous {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Could unmarshal value to map[string]tftypes.Value",
			Detail:   "Could unmarshal value to map[string]tftypes.Value to extract guid_seed_addition",
		})
		goto Return
	}
	_ = value

	if seedAdditionValue, isSuccesful = valueMap["guid_seed_addition"]; !isSuccesful {
		// Not having guid_seed_addition is fine.
		goto Return
	}
	_ = seedAdditionValue

	if !seedAdditionValue.IsKnown() {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Provider configuration has 'guid_seed_addition' attribute but it is not known at plan-time.",
			Detail:   "The 'guid_seed_addition' attribute must be known during the plan phase. See attribute description for more informations.",
		})

		goto Return
	}

	if seedAdditionValue.IsNull() {
		goto Return
	}

	if err := seedAdditionValue.As(&seedAddition); err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Extraction failed",
			Detail:   "The guid seed addition could not be extracted as string",
		})

		goto Return
	} else {
		workdir, _ := os.Getwd()
		seedAddition = strings.ReplaceAll(seedAddition, "{workdir}", workdir)
	}

Return:
	return seedAddition, diags, len(diags) == 0
}

func GetGuidSeedAdditionAttributeDescription() string {
	return "It serves as an guid seed addition to those resources that implement `guid_seed` as an " +
		"attribute. But there are scopes you need to keep in mind: if `guid_seed_addition` has been " +
		"specified in the provider block then top-level and nested modules are using the provider " +
		"block seed addition. If `guid_seed_addition` has been specified in the provider_meta block " +
		"then only the resources of that module are using the module-level seed addition. " +
		"Besides `guid_seed`, the provider block seed addition, the provider_meta block seed addition " +
		"and the resource type itself will become part of the final seed.\n" + `
	**Placeholders**:
	- "{workdir}" (Keyword) The actual workdir; equals to terraform's path.root. This placeholder is
	recommended because this value won't be dragged along the plan and apply phase in comparison to
	"abspath(path.root)" that you would add to resource seed where a change to path.root would be
	recognized just as usual from terraform.`
}

func GetGuidSeedAdditionSchemaAttribute(attributeDescription string) *tfprotov6.SchemaAttribute {
	return &tfprotov6.SchemaAttribute{
		Name:        "guid_seed_addition",
		Type:        tftypes.String,
		Required:    false,
		Optional:    true,
		Computed:    false,
		Description: attributeDescription,
	}
}
