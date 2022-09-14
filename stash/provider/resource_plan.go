package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// PlanResourceChange function
func (s *RawProviderServer) PlanResourceChange(ctx context.Context, request *tfprotov6.PlanResourceChangeRequest) (*tfprotov6.PlanResourceChangeResponse, error) {
	response := &tfprotov6.PlanResourceChangeResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		response.Diagnostics = append(response.Diagnostics, execDiag...)
		return response, nil
	}

	_, err := GetResourceType(request.TypeName)

	if err != nil {
		response.Diagnostics = append(response.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to determine planned resource type",
			Detail:   err.Error(),
		})

		return response, nil
	}

	response.PlannedState = request.ProposedNewState
	return response, nil
}
