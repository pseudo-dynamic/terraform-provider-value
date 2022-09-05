package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ReadResource function
func (s *RawProviderServer) ReadResource(ctx context.Context, req *tfprotov5.ReadResourceRequest) (*tfprotov5.ReadResourceResponse, error) {
	resp := &tfprotov5.ReadResourceResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		resp.Diagnostics = append(resp.Diagnostics, execDiag...)
		return resp, nil
	}

	var resState map[string]tftypes.Value
	var err error
	rt, err := GetResourceType(req.TypeName)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to determine resource type",
			Detail:   err.Error(),
		})
		
		return resp, nil
	}

	currentState, err := req.CurrentState.Unmarshal(rt)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to decode current state",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	if currentState.IsNull() {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to read resource",
			Detail:   "Incomplete or missing state",
		})

		return resp, nil
	}

	err = currentState.As(&resState)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to extract resource from current state",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	_, isValueExisting := resState["value"]

	if !isValueExisting {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Current state of resource has no 'value' attribute",
			Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, plese report issue to provider maintainers.",
		})

		return resp, nil
	}

	_, isResultExisting := resState["result"]

	if !isResultExisting {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Current state of resource has no 'result' attribute",
			Detail:   "This should not happen. The state may be incomplete or corrupted.\nIf this error is reproducible, plese report issue to provider maintainers.",
		})

		return resp, nil
	}

	resp.NewState = req.CurrentState
	return resp, nil
}
