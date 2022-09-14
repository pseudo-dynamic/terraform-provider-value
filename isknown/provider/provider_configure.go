package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"golang.org/x/mod/semver"
)

const minTFVersion string = "v0.14.8"

// ConfigureProvider function
func (s *UserProviderServer) ConfigureProvider(ctx context.Context, req *tfprotov6.ConfigureProviderRequest) (*tfprotov6.ConfigureProviderResponse, error) {
	return &tfprotov6.ConfigureProviderResponse{}, nil
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
