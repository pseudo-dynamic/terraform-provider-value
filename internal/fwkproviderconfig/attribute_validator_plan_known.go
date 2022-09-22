package fwkproviderconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type PlanKnownValidator struct{}

func (v *PlanKnownValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig != nil && req.AttributeConfig.IsUnknown() {
		resp.Diagnostics.AddError(
			req.AttributePath.String()+" is unknown",
			req.AttributePath.String()+" must be fully known at plan-time. For further informations take a look into the documentation.")
	}
}

func getPlanKnownValidatorDescription() string {
	return "Validates that the value is not unknown."
}

func (*PlanKnownValidator) Description(context.Context) string {
	return getPlanKnownValidatorDescription()
}

func (*PlanKnownValidator) MarkdownDescription(context.Context) string {
	return getPlanKnownValidatorDescription()
}
