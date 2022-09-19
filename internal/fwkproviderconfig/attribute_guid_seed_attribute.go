package fwkproviderconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type guidSeedAdditionUnknownValidator struct{}

func (v *guidSeedAdditionUnknownValidator) Validate(ctx context.Context, req tfsdk.ValidateAttributeRequest, resp *tfsdk.ValidateAttributeResponse) {
	if req.AttributeConfig != nil && req.AttributeConfig.IsUnknown() {
		resp.Diagnostics.AddError(
			"Guid seed addition is unknown",
			"Guid seed addition must be fully known at plan-time. For further informations take a look into the documentation.")
	}
}

func getGuidSeedAdditionUnknownValidatorDescription() string {
	return "Validates that guid seed addition is not unknown."
}

func (*guidSeedAdditionUnknownValidator) Description(context.Context) string {
	return getGuidSeedAdditionUnknownValidatorDescription()
}

func (*guidSeedAdditionUnknownValidator) MarkdownDescription(context.Context) string {
	return getGuidSeedAdditionUnknownValidatorDescription()
}

type guidSeedAdditionDefaultEmptyModifier struct{}

func (r guidSeedAdditionDefaultEmptyModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
	if req.AttributePlan == nil || req.AttributePlan.IsNull() {
		resp.AttributePlan = types.String{Value: ""}
	}
}

func getGuidSeedAdditionAttributeModifierDescription() string {
	return "If the value is null, then its value is now empty."
}

// Description returns a human-readable description of the plan modifier.
func (guidSeedAdditionDefaultEmptyModifier) Description(context.Context) string {
	return getGuidSeedAdditionAttributeModifierDescription()
}

// MarkdownDescription returns a markdown description of the plan modifier.
func (guidSeedAdditionDefaultEmptyModifier) MarkdownDescription(context.Context) string {
	return getGuidSeedAdditionAttributeModifierDescription()
}
