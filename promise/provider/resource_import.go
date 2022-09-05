package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// ImportResourceState function
func (s *RawProviderServer) ImportResourceState(ctx context.Context, req *tfprotov5.ImportResourceStateRequest) (*tfprotov5.ImportResourceStateResponse, error) {
	// Terraform only gives us the schema name of the resource and an ID string, as passed by the user on the command line.
	// The ID should be a combination of a Kubernetes GVK and a namespace/name type of resource identifier.
	// Without the user supplying the GRV there is no way to fully identify the resource when making the Get API call to K8s.
	// Presumably the Kubernetes API machinery already has a standard for expressing such a group. We should look there first.
	resp := &tfprotov5.ImportResourceStateResponse{}

	resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
		Severity: tfprotov5.DiagnosticSeverityError,
		Summary:  "Import not supported",
	})

	return resp, nil
}
