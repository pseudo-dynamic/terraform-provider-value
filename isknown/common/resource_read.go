package common

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// CanReadResource function
func CanReadResource(
	currentState *tfprotov6.DynamicValue,
	resourceType tftypes.Type,
	resp *tfprotov6.ReadResourceResponse) (
	tftypes.Value,
	map[string]tftypes.Value,
	bool) {
	var currentStateValue tftypes.Value
	var currentStateValueMap map[string]tftypes.Value
	var err error
	currentStateValue, err = currentState.Unmarshal(resourceType)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to decode current state",
			Detail:   err.Error(),
		})

		return currentStateValue, currentStateValueMap, false
	}

	if currentStateValue.IsNull() {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to read resource",
			Detail:   "Incomplete or missing state",
		})

		return currentStateValue, currentStateValueMap, false
	}

	err = currentStateValue.As(&currentStateValueMap)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to extract resource from current state",
			Detail:   err.Error(),
		})

		return currentStateValue, currentStateValueMap, false
	}

	_, isValueExisting := currentStateValueMap["value"]

	if !isValueExisting {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Current state of resource has no 'value' attribute",
			Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, please report issue to provider maintainers.",
		})

		return currentStateValue, currentStateValueMap, false
	}

	_, isResultExisting := currentStateValueMap["result"]

	if !isResultExisting {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Current state of resource has no 'result' attribute",
			Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, please report issue to provider maintainers.",
		})

		return currentStateValue, currentStateValueMap, false
	}

	_, isSeedExisting := currentStateValueMap["seed"]

	if !isSeedExisting {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Current state of resource has no 'seed' attribute",
			Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, please report issue to provider maintainers.",
		})

		return currentStateValue, currentStateValueMap, false
	}

	// if !currentStateValueMap["seed"].IsFullyKnown() {
	// 	resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
	// 		Severity: tfprotov6.DiagnosticSeverityError,
	// 		Summary:  "Current state of resource has a 'seed' attribute but it is not fully known.",
	// 		Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, please report issue to provider maintainers.",
	// 	})

	// 	return currentStateValue, currentStateValueMap, false
	// }

	return currentStateValue, currentStateValueMap, true
}
