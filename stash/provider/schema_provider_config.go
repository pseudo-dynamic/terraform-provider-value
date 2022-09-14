package provider

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
)

// GetProviderConfigSchema contains the definitions of all configuration attributes
func GetProviderConfigSchema() *tfprotov6.Schema {
	b := tfprotov6.SchemaBlock{}

	return &tfprotov6.Schema{
		Version: 0,
		Block:   &b,
	}
}
