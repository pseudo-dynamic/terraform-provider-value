package provider

import (
	"context"
	// "fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	// "github.com/hashicorp/terraform-plugin-go/tftypes"
)

// PlanResourceChange function
func (s *RawProviderServer) PlanResourceChange(ctx context.Context, request *tfprotov5.PlanResourceChangeRequest) (*tfprotov5.PlanResourceChangeResponse, error) {
	response := &tfprotov5.PlanResourceChangeResponse{}
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

	response.PlannedState = request.ProposedNewState
	return response, nil

	// // Decode proposed resource state
	// proposedState, err := request.ProposedNewState.Unmarshal(resourceType)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal planned resource state",
	// 		Detail:   err.Error(),
	// 	})

	// 	return response, nil
	// }

	// proposedValue := make(map[string]tftypes.Value)
	// err = proposedState.As(&proposedValue)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to extract planned resource state from tftypes.Value",
	// 		Detail:   err.Error(),
	// 	})

	// 	return response, nil
	// }

	// // Decode prior resource state
	// priorState, err := request.PriorState.Unmarshal(resourceType)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to unmarshal prior resource state",
	// 		Detail:   err.Error(),
	// 	})

	// 	return response, nil
	// }

	// s.logger.Trace("[PlanResourceChange]", "[PriorState]", dump(priorState))
	// priorValue := make(map[string]tftypes.Value)
	// err = priorState.As(&priorValue)

	// if err != nil {
	// 	response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 		Severity: tfprotov5.DiagnosticSeverityError,
	// 		Summary:  "Failed to extract prior resource state from tftypes.Value",
	// 		Detail:   err.Error(),
	// 	})

	// 	return response, nil
	// }

	// if proposedState.IsNull() {
	// 	// we plan to delete the resource
	// 	if _, ok := priorValue["timestamp"]; ok {
	// 		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 			Severity: tfprotov5.DiagnosticSeverityError,
	// 			Summary:  "Invalid prior state while planning for destroy",
	// 			Detail:   fmt.Sprintf("'timestamp' attribute missing from state: %s", err),
	// 		})

	// 		return response, nil
	// 	}

	// 	response.PlannedState = request.ProposedNewState
	// 	return response, nil
	// }

	// if proposedValue["timestamp"].IsNull() {
	// 	// plan for Create
	// 	proposedValue["timestamp"] = tftypes.NewValue(tftypes.String, tftypes.UnknownValue)
	// 	propStateVal := tftypes.NewValue(proposedState.Type(), proposedValue)
	// 	s.logger.Trace("[PlanResourceChange]", "new planned state", dump(propStateVal))
	// 	plannedState, err := tfprotov5.NewDynamicValue(resourceType, propStateVal)

	// 	if err != nil {
	// 		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 			Severity: tfprotov5.DiagnosticSeverityError,
	// 			Summary:  "Failed to assemble proposed state during plan",
	// 			Detail:   err.Error(),
	// 		})

	// 		return response, nil
	// 	}

	// 	response.PlannedState = &plannedState
	// } else {
	// 	// plan for Update
	// 	// NO-OP
	// 	response.PlannedState = request.PriorState
	// }

	// return response, nil
}
