package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ValidateDataResourceConfig function
func (s *UserProviderServer) ValidateDataResourceConfig(ctx context.Context, req *tfprotov6.ValidateDataResourceConfigRequest) (*tfprotov6.ValidateDataResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateDataResourceConfigResponse{}
	return resp, nil
}

// ValidateResourceConfig function
func (s *UserProviderServer) ValidateResourceConfig(ctx context.Context, req *tfprotov6.ValidateResourceConfigRequest) (*tfprotov6.ValidateResourceConfigResponse, error) {
	resp := &tfprotov6.ValidateResourceConfigResponse{}
	resourceType := getResourceType(req.TypeName)

	// Decode proposed resource state
	config, err := req.Config.Unmarshal(resourceType)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to unmarshal resource state",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	att := tftypes.NewAttributePath()
	att = att.WithAttributeName("value")

	configVal := make(map[string]tftypes.Value)
	err = config.As(&configVal)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to extract resource state from SDK value",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	_, ok := configVal["value"]

	if !ok {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Value missing from resource configuration",
			Detail:    "A value attribute containing a valid terraform value is required.",
			Attribute: att,
		})

		return resp, nil
	}

	_, ok = configVal["guid_seed"]

	if !ok {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity:  tfprotov6.DiagnosticSeverityError,
			Summary:   "Unique seed missing from resource configuration",
			Detail:    "A guid_seed attribute containing a valid terraform value is required.",
			Attribute: att,
		})

		return resp, nil
	}

	return resp, nil
}
