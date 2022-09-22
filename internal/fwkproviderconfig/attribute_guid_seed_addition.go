package fwkproviderconfig

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
)

const GuidSeedAdditionAttributeName string = "guid_seed_addition"

func getGuidSeedAdditionAttribute() tfsdk.Attribute {
	return tfsdk.Attribute{
		Type:     types.StringType,
		Required: false,
		Optional: true,
		Computed: false,
		Validators: []tfsdk.AttributeValidator{
			&PlanKnownValidator{},
		},
		PlanModifiers: tfsdk.AttributePlanModifiers{
			guidSeedAdditionDefaultEmptyModifier{},
		},
		Description: goproviderconfig.GetGuidSeedAdditionAttributeDescription(),
	}
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
