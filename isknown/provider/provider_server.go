package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// ValidateProviderConfig function
func (s *UserProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	resp := &tfprotov6.ValidateProviderConfigResponse{PreparedConfig: req.Config}
	return resp, nil
}

// UpgradeResourceState isn't really useful in this provider, but we have to loop the state back through to keep Terraform happy.
func (s *UserProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov6.UpgradeResourceStateRequest) (*tfprotov6.UpgradeResourceStateResponse, error) {
	resp := &tfprotov6.UpgradeResourceStateResponse{}
	resp.Diagnostics = []*tfprotov6.Diagnostic{}

	sch := getProviderResourceSchema(req.TypeName)
	rt := getObjectTypeFromSchema(sch[req.TypeName])

	rv, err := req.RawState.Unmarshal(rt)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to decode old state during upgrade",
			Detail:   err.Error(),
		})
		return resp, nil
	}

	us, err := tfprotov6.NewDynamicValue(rt, rv)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Failed to encode new state during upgrade",
			Detail:   err.Error(),
		})
	}

	resp.UpgradedState = &us
	return resp, nil
}

// ReadDataSource function
func (s *UserProviderServer) ReadDataSource(ctx context.Context, req *tfprotov6.ReadDataSourceRequest) (*tfprotov6.ReadDataSourceResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReadDataSource not implemented")
}

// StopProvider function
func (s *UserProviderServer) StopProvider(ctx context.Context, req *tfprotov6.StopProviderRequest) (*tfprotov6.StopProviderResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method StopProvider not implemented")
}
