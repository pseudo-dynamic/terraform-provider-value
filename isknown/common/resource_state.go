package common

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

func UnmarshalState(state *tfprotov6.DynamicValue, stateType tftypes.Type) (tftypes.Value, map[string]tftypes.Value, []*tfprotov6.Diagnostic, bool) {
	diags := []*tfprotov6.Diagnostic{}
	resourceValueMap := make(map[string]tftypes.Value)
	resourceValue, err := state.Unmarshal(stateType)

	if err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to unmarshal state",
			Detail:   err.Error(),
		})

		goto End
	}

	err = resourceValue.As(&resourceValueMap)

	if err != nil {
		diags = append(diags, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to extract state from tftypes.Value",
			Detail:   err.Error(),
		})

		goto End
	}

End:
	return resourceValue, resourceValueMap, diags, len(diags) != 0
}
