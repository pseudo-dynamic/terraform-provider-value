package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
	"golang.org/x/mod/semver"
)

const minTFVersion string = "v0.14.8"

// ConfigureProvider function
func (s *RawProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov5.ConfigureProviderRequest) (*tfprotov5.ConfigureProviderResponse, error) {
	return &tfprotov5.ConfigureProviderResponse{}, nil
}

func (s *RawProviderServer) canExecute() (resp []*tfprotov5.Diagnostic) {
	if semver.IsValid(s.hostTFVersion) && semver.Compare(s.hostTFVersion, minTFVersion) < 0 {
		resp = append(resp, &tfprotov5.Diagnostic{
			Severity: tfprotov5.DiagnosticSeverityError,
			Summary:  "Incompatible terraform version",
			Detail:   fmt.Sprintf("The `is_known` resource requires Terraform %s or above", minTFVersion),
		})
	}

	return
}
