package fwkproviderconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
)

func GetProviderConfigSchema() *tfsdk.Schema {
	return &tfsdk.Schema{
		Version: 0,
		Attributes: map[string]tfsdk.Attribute{
			"guid_seed_addition": {
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
			},
		},
	}
}
