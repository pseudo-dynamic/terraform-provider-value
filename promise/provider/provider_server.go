package provider

import (
	"context"

	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// RawProviderServer implements the ProviderServer interface as exported from ProtoBuf.
type RawProviderServer struct {
	// Since the provider is essentially a gRPC server, the execution flow is dictated by the order of the client (Terraform) request calls.
	// Thus it needs a way to persist state between the gRPC calls. These attributes store values that need to be persisted between gRPC calls,
	// such as instances of the Kubernetes clients, configuration options needed at runtime.
	logger hclog.Logger

	//providerEnabled bool
	hostTFVersion string
}

func dump(v interface{}) hclog.Format {
	return hclog.Fmt("%v", v)
}

// PrepareProviderConfig function
func (s *RawProviderServer) PrepareProviderConfig(ctx context.Context, req *tfprotov5.PrepareProviderConfigRequest) (*tfprotov5.PrepareProviderConfigResponse, error) {
	s.logger.Trace("[PrepareProviderConfig][Request]\n%s\n", dump(*req))
	resp := &tfprotov5.PrepareProviderConfigResponse{PreparedConfig: req.Config}
	return resp, nil
}

// ValidateDataSourceConfig function
func (s *RawProviderServer) ValidateDataSourceConfig(ctx context.Context, req *tfprotov5.ValidateDataSourceConfigRequest) (*tfprotov5.ValidateDataSourceConfigResponse, error) {
	s.logger.Trace("[ValidateDataSourceConfig][Request]\n%s\n", dump(*req))
	resp := &tfprotov5.ValidateDataSourceConfigResponse{}
	return resp, nil
}

// UpgradeResourceState isn't really useful in this provider, but we have to loop the state back through to keep Terraform happy.
func (s *RawProviderServer) UpgradeResourceState(ctx context.Context, req *tfprotov5.UpgradeResourceStateRequest) (*tfprotov5.UpgradeResourceStateResponse, error) {
	resp := &tfprotov5.UpgradeResourceStateResponse{}
	resp.Diagnostics = []*tfprotov5.Diagnostic{}

	sch := GetProviderResourceSchema()
	rt := GetObjectTypeFromSchema(sch[req.TypeName])

	rv, err := req.RawState.Unmarshal(rt)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to decode old state during upgrade",
			Detail:   err.Error(),
		})
		return resp, nil
	}

	us, err := tfprotov5.NewDynamicValue(rt, rv)

	if err != nil {
		resp.Diagnostics = append(resp.Diagnostics, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Failed to encode new state during upgrade",
			Detail:   err.Error(),
		})
	}

	resp.UpgradedState = &us
	return resp, nil
}

// ReadDataSource function
func (s *RawProviderServer) ReadDataSource(ctx context.Context, req *tfprotov5.ReadDataSourceRequest) (*tfprotov5.ReadDataSourceResponse, error) {
	s.logger.Trace("[ReadDataSource][Request]\n%s\n", dump(*req))
	return nil, status.Errorf(codes.Unimplemented, "method ReadDataSource not implemented")
}

// StopProvider function
func (s *RawProviderServer) StopProvider(ctx context.Context, req *tfprotov5.StopProviderRequest) (*tfprotov5.StopProviderResponse, error) {
	s.logger.Trace("[StopProvider][Request]\n%s\n", dump(*req))
	return nil, status.Errorf(codes.Unimplemented, "method StopProvider not implemented")
}
