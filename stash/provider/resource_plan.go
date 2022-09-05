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

	// proposedValueMap := make(map[string]tftypes.Value)
	// err = proposedState.As(&proposedValueMap)

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
	// 	// Plan to delete
	// 	response.PlannedState = request.ProposedNewState
	// 	return response, nil
	// }

	// if proposedValueMap["value"].IsNull() {
	// 	proposedValueMap["result"] = tftypes.NewValue(tftypes.DynamicPseudoType, nil)
	// }

	// if proposedValueMap["value"].IsNull() {
	// 	proposedValueMap["value"] = tftypes.NewValue(tftypes.DynamicPseudoType)
	// 	customProposedValue := tftypes.NewValue(proposedState.Type(), proposedValueMap)
	// 	s.logger.Trace("[PlanResourceChange]", "new planned state", dump(customProposedValue))
	// 	customPlannedState, err := tfprotov5.NewDynamicValue(resourceType, customProposedValue)

	// 	if err != nil {
	// 		response.Diagnostics = append(response.Diagnostics, &tfprotov5.Diagnostic{
	// 			Severity: tfprotov5.DiagnosticSeverityError,
	// 			Summary:  "Failed to assemble proposed state during plan",
	// 			Detail:   err.Error(),
	// 		})

	// 		return response, nil
	// 	}

	// 	response.PlannedState = &customPlannedState
	// } else {
	// 	response.PlannedState = request.ProposedNewState
	// }

	// customProposedValue := tftypes.NewValue(proposedState.Type(), proposedValueMap)
	// s.logger.Trace("[PlanResourceChange]", "new planned state", dump(customProposedValue))
	// customPlannedState, err := tfprotov5.NewDynamicValue(resourceType, customProposedValue)
	// response.PlannedState = &customPlannedState
	// return response, nil
}
