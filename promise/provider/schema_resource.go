package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"
)

// GetResourceType returns the tftypes.Type of a resource of type 'name'
func GetResourceType(name string) (tftypes.Type, error) {
	sch := GetProviderResourceSchema()
	rsch, ok := sch[name]

	if !ok {
		return tftypes.DynamicPseudoType, fmt.Errorf("unknown resource %s - cannot find schema", name)
	}

	return schema.GetObjectTypeFromSchema(rsch), nil
}

// GetProviderResourceSchema contains the definitions of all supported resources
func GetProviderResourceSchema() map[string]*tfprotov6.Schema {
	return map[string]*tfprotov6.Schema{
		"value_promise": {
			Version: 1,
			Block: &tfprotov6.SchemaBlock{
				Description: "Allows you to treat a value as unknown. This is desirable when you want postconditions first being evaluated during apply phase.",
				BlockTypes:  []*tfprotov6.SchemaNestedBlock{},
				Attributes: []*tfprotov6.SchemaAttribute{
					{
						Name:        "value",
						Type:        tftypes.DynamicPseudoType,
						Required:    true,
						Optional:    false,
						Computed:    false,
						Description: "The value to promise. Any (nested) change to `value` results into `result` to be marked as \"(known after apply)\"",
					},
					{
						Name:        "result",
						Type:        tftypes.DynamicPseudoType,
						Required:    false,
						Optional:    false,
						Computed:    true,
						Description: "`result` is as soon as you apply set to `value`. Every change of `value` results into `result` to be marked as \"(known after apply)\"",
					},
				},
			},
		},
	}
}
