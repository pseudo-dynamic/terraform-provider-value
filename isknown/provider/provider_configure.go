package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
	"golang.org/x/mod/semver"
)

const minTFVersion string = "v0.14.8"

// ValidateProviderConfig function
func (s *UserProviderServer) ValidateProviderConfig(ctx context.Context, req *tfprotov6.ValidateProviderConfigRequest) (*tfprotov6.ValidateProviderConfigResponse, error) {
	resp := &tfprotov6.ValidateProviderConfigResponse{PreparedConfig: req.Config}
	return resp, nil
}

// ConfigureProvider function
func (s *UserProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	var isWorking bool
	resp := &tfprotov6.ConfigureProviderResponse{}
	var diags []*tfprotov6.Diagnostic

	var providerConfigSeedAddition string
	if providerConfigSeedAddition, diags, isWorking = goproviderconfig.TryExtractProviderConfigGuidSeedAddition(req.Config); !isWorking {
		resp.Diagnostics = append(resp.Diagnostics, diags...)
		return resp, nil
	}

	s.ProviderConfigSeedAddition = providerConfigSeedAddition
	return resp, nil
}

func (s *UserProviderServer) canExecute() (resp []*tfprotov6.Diagnostic) {
	if semver.IsValid(s.hostTFVersion) && semver.Compare(s.hostTFVersion, minTFVersion) < 0 {
		resp = append(resp, &tfprotov6.Diagnostic{
			Severity: tfprotov6.DiagnosticSeverityError,
			Summary:  "Incompatible terraform version",
			Detail:   fmt.Sprintf("The resource requires Terraform %s or above", minTFVersion),
		})
	}

	return
}
