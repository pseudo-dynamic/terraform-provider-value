package provider

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov5"
)

// GetProviderConfigSchema contains the definitions of all configuration attributes
func GetProviderConfigSchema() *tfprotov5.Schema {
	b := tfprotov5.SchemaBlock{}

	return &tfprotov5.Schema{
		Version: 0,
		Block:   &b,
	}
}
