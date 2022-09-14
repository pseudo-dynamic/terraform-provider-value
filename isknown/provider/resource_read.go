package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// ReadResource function
func (s *UserProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		resp.Diagnostics = append(resp.Diagnostics, execDiag...)
		return resp, nil
	}

	resp.NewState = req.CurrentState
	return resp, nil
}
