package provider

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// ValidateResourceTypeConfig function
func (s *RawProviderServer) ValidateResourceTypeConfig(ctx context.Context, req *tfprotov5.ValidateResourceTypeConfigRequest) (*tfprotov5.ValidateResourceTypeConfigResponse, error) {
	resp := &tfprotov5.ValidateResourceTypeConfigResponse{}
	rt, err := GetResourceType(req.TypeName)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to determine resource type",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	log.Println("--------------------TEST----------------------")
	log.Printf("ResourceType: %v\n", rt)

	// Decode proposed resource state
	config, err := req.Config.Unmarshal(rt)
	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
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
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to extract resource state from SDK value",
			Detail:   err.Error(),
		})

		return resp, nil
	}

	_, ok := configVal["value"]

	if !ok {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity:  tfprotov5.DiagnosticSeverityError,
			Summary:   "Value missing from resource configuration",
			Detail:    "A value attribute containing a valid terraform value is required.",
			Attribute: att,
		})

		return resp, nil
	}

	return resp, nil
}
