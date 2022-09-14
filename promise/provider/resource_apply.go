package provider

import (
	"context"
	// "fmt"
	// "time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ApplyResourceChange function
func (s *RawProviderServer) ApplyResourceChange(ctx context.Context, req *tfprotov6.ApplyResourceChangeRequest) (*tfprotov6.ApplyResourceChangeResponse, error) {
	response := &tfprotov6.ApplyResourceChangeResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		response.Diagnostics = append(response.Diagnostics, execDiag...)
		return response, nil
	}

	rt, err := GetResourceType(req.TypeName)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to determine planned resource type",
			Detail:   err.Error(),
		})

		return response, nil
	}

	// response.NewState = request.PlannedState
	// return response, nil

	plannedState, err := req.PlannedState.Unmarshal(rt)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to unmarshal planned resource state",
			Detail:   err.Error(),
		})
		return response, nil
	}

	priorState, err := req.PriorState.Unmarshal(rt)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to unmarshal prior resource state",
			Detail:   err.Error(),
		})
		return response, nil
	}

	plannedValueMap := make(map[string]tftypes.Value)
	err = plannedState.As(&plannedValueMap)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to extract planned resource state from tftypes.Value",
			Detail:   err.Error(),
		})
		return response, nil
	}

	switch {
	case priorState.IsNull():
		// This is a "create"
		fallthrough
	case !plannedState.IsNull() && !priorState.IsNull():
		// This is a "create" OR "update"
		plannedValueMap["result"] = plannedValueMap["value"]
		customPlannedValue := tftypes.NewValue(plannedState.Type(), plannedValueMap)
		customPlannedState, err := tfprotov6.NewDynamicValue(rt, customPlannedValue)

		if err != nil {
			response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
				Severity: tfprotov6.DiagnosticSeverityError,
				Summary:  "Failed to assemble proposed state during apply/update",
				Detail:   err.Error(),
			})
			return response, nil
		}

		response.NewState = &customPlannedState
		return response, nil
	case plannedState.IsNull():
		// Delete the resource
		break
	}

	return response, nil
}
