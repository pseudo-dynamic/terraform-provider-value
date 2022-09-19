package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// ReadResource function
func (s *UserProviderServer) ReadResource(ctx context.Context, req *tfprotov6.ReadResourceRequest) (*tfprotov6.ReadResourceResponse, error) {
	resp := &tfprotov6.ReadResourceResponse{}
	execDiag := s.canExecute()

	if len(execDiag) > 0 {
		resp.Diagnostics = append(resp.Diagnostics, execDiag...)
		return resp, nil
	}

	resourceType := getResourceType(req.TypeName)

	if _, _, canReadCurrentState := common.TryReadResource(req.CurrentState, resourceType, resp); !canReadCurrentState {
		return resp, nil
	}

	resp.NewState = req.CurrentState
	return resp, nil
}
