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

	// plannedState, err := request.PlannedState.Unmarshal(resourceType)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal planned resource state",
	// 		Detail:   err.Error(),
	// 	})
	// 	return response, nil
	// }

	// s.logger.Trace("[ApplyResourceChange][PlannedState] %#v", plannedState)
	// priorState, err := request.PriorState.Unmarshal(resourceType)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal prior resource state",
	// 		Detail:   err.Error(),
	// 	})
	// 	return response, nil
	// }

	// s.logger.Trace("[ApplyResourceChange]", "[PriorState]", dump(priorState))
	// plannedValueMap := make(map[string]tftypes.Value)
	// err = plannedState.As(&plannedValueMap)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to extract planned resource state from tftypes.Value",
	// 		Detail:   err.Error(),
	// 	})
	// 	return response, nil
	// }

	// if !plannedState.IsNull() {
	// 	// Create or update
	// 	if plannedValueMap["value"].IsNull() {
	// 		plannedValueMap["value"] = tftypes.NewValue(tftypes.DynamicPseudoType, nil)
	// 	}
	// }

	// switch {
	// case applyPriorState.IsNull():
	// 	// This is a "create"
	// 	// All we need to do is update the timestamp
	// 	timestamp := time.Now().Unix()
	// 	applyPlannedValue["timestamp"] = tftypes.NewValue(tftypes.String, fmt.Sprint(timestamp))

	// 	applyStateVal := tftypes.NewValue(applyPlannedState.Type(), applyPlannedValue)
	// 	s.logger.Trace("[ApplyResourceChange]", "[PropStateVal]", dump(applyStateVal))
	// 	customPlannedState, err := tfprotov5.NewDynamicValue(resourceType, applyStateVal)

	// 	if err != nil {
	// 		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 			Severity: tfprotov5.DiagnosticSeverityError,
	// 			Summary:  "Failed to assemble proposed state during apply",
	// 			Detail:   err.Error(),
	// 		})
	// 		return response, nil
	// 	}

	// 	s.logger.Trace("[ApplyResourceChange]", "[PlannedState]", dump(customPlannedState))
	// 	response.NewState = &customPlannedState
	// case !applyPlannedState.IsNull() && !applyPriorState.IsNull():
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Attempting to perform update on cache resource",
	// 		Detail:   "An update operation was attempted on a cache resource. This should not occur. Please report this to the provider maintainers.",
	// 	})

	// 	return response, nil
	// case applyPlannedState.IsNull():
	// 	// Delete the resource
	// 	return response, nil
	// }

	// return response, nil
}
