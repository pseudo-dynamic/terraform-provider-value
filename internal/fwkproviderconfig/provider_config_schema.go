package fwkproviderconfig

import (
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

func GetProviderConfigSchema() *tfsdk.Schema {
	return &tfsdk.Schema{
		Version: 0,
		Attributes: map[string]tfsdk.Attribute{
			GuidSeedAdditionAttributeName: getGuidSeedAdditionAttribute(),
		},
	}
}
