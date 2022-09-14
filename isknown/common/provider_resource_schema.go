package common

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

type ProviderResourceSchemaParameters struct {
	ResourceName        string
	ResourceDescription string
	ValueDescription    string
	ResultDescription   string
}

// GetProviderResourceSchema contains the definitions of all supported resources
func GetProviderResourceSchema(schema ProviderResourceSchemaParameters) map[string]*tfprotov6.Schema {
	return map[string]*tfprotov6.Schema{
		schema.ResourceName: {
			Version: 1,
			Block: &tfprotov6.SchemaBlock{
				Description: schema.ResourceDescription + "\n" + seedPrefixDescription,
				BlockTypes:  []*tfprotov6.SchemaNestedBlock{},
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name:        "value",
						Type:        tftypes.DynamicPseudoType,
						Required:    true,
						Optional:    false,
						Computed:    false,
						Description: schema.ValueDescription,
					},
					{
						Name:     "unique_seed",
						Type:     tftypes.String,
						Required: true,
						Optional: false,
						Computed: false,
						Description: "Attention! The seed is being used to determine resource uniqueness prior and " +
							"during apply-phase. Very important to state is that the **seed must be fully known during " +
							"the plan phase**, otherwise, an error is thrown. Within one terraform plan & apply the **seed " +
							"of every \"" + schema.ResourceName + "\" must be unique**! I recommend you to use the " +
							"provider_meta-feature for increased uniqueness. Under certain circumstances you may " +
							"face problems if you run terraform concurrenctly. If you do so, " +
							"then I recommend you to pass-through a random value via a user (environment) variable " +
							"that you then add to the seed.",
					},
					{
						Name:     "proposed_unknown",
						Type:     tftypes.DynamicPseudoType,
						Required: true,
						Optional: false,
						Computed: false,
						Description: "It is very crucial that this field is **not** filled by any " +
							"custom value except the one produced by `value_unknown_proposer` (resource). " +
							"This has the reason as its `value` is **always** unknown during the plan phase. " +
							"On this behaviour this resource must rely and it cannot check if you do not so!",
					},
					{
						Name:        "result",
						Type:        tftypes.Bool,
						Required:    false,
						Optional:    true,
						Computed:    true,
						Description: schema.ResultDescription,
					},
				},
			},
		},
	}
}
