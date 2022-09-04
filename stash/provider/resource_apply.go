package provider

import (
	"context"
	// "fmt"
	// "time"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	// "github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ApplyResourceChange function
func (s *RawProviderServer) ApplyResourceChange(ctx context.Context, request *tfprotov5.ApplyResourceChangeRequest) (*tfprotov5.ApplyResourceChangeResponse, error) {
	response := &tfprotov5.ApplyResourceChangeResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		response.Diagnostics = append(response.Diagnostics, execDiag...)
		return response, nil
	}

	_, err := GetResourceType(request.TypeName)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to determine planned resource type",
			Detail:   err.Error(),
		})

		return response, nil
	}

	response.NewState = request.PlannedState
	return response, nil

	// applyPlannedState, err := req.PlannedState.Unmarshal(rt)

	// if err != nil {
	// 	resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal planned resource state",
	// 		Detail:   err.Error(),
	// 	})
	// 	return resp, nil
	// }

	// s.logger.Trace("[ApplyResourceChange][PlannedState] %#v", applyPlannedState)

	// applyPriorState, err := req.PriorState.Unmarshal(rt)

	// if err != nil {
	// 	resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal prior resource state",
	// 		Detail:   err.Error(),
	// 	})
	// 	return resp, nil
	// }

	// s.logger.Trace("[ApplyResourceChange]", "[PriorState]", dump(applyPriorState))

	// applyPlannedValue := make(map[string]tftypes.Value)
	// err = applyPlannedState.As(&applyPlannedValue)

	// if err != nil {
	// 	resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to extract planned resource state from tftypes.Value",
	// 		Detail:   err.Error(),
	// 	})
	// 	return resp, nil
	// }

	// switch {
	// case applyPriorState.IsNull():
	// 	// This is a "create"
	// 	// All we need to do is update the timestamp
	// 	timestamp := time.Now().Unix()
	// 	applyPlannedValue["timestamp"] = tftypes.NewValue(tftypes.String, fmt.Sprint(timestamp))

	// 	applyStateVal := tftypes.NewValue(applyPlannedState.Type(), applyPlannedValue)
	// 	s.logger.Trace("[ApplyResourceChange]", "[PropStateVal]", dump(applyStateVal))
	// 	plannedState, err := tfprotov5.NewDynamicValue(rt, applyStateVal)

	// 	if err != nil {
	// 		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
	// 			Severity: tfprotov5.DiagnosticSeverityError,
	// 			Summary:  "Failed to assemble proposed state during apply",
	// 			Detail:   err.Error(),
	// 		})
	// 		return resp, nil
	// 	}

	// 	s.logger.Trace("[ApplyResourceChange]", "[PlannedState]", dump(plannedState))
	// 	resp.NewState = &plannedState
	// case !applyPlannedState.IsNull() && !applyPriorState.IsNull():
	// 	resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Attempting to perform update on cache resource",
	// 		Detail:   "An update operation was attempted on a cache resource. This should not occur. Please report this to the provider maintainers.",
	// 	})

	// 	return resp, nil
	// case applyPlannedState.IsNull():
	// 	// Delete the resource
	// 	return resp, nil
	// }

	// return resp, nil
}
