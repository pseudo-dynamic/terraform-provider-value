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
		goto Return
	}
	_ = value

	if seedAdditionValue, isSuccesful = valueMap["guid_seed_addition"]; !isSuccesful {
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
	return "It serves as addition to each seed of any `value_is_fully_known` (resource) or " +
		"`value_is_known` (resource) within the project if specified in provider, or within the " +
		"same module if specified in provider-meta.\n" + `
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
