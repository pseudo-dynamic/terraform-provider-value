package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
)

// GetObjectTypeFromSchema returns a tftypes.Type that can wholy represent the schema input
func GetObjectTypeFromSchema(schema *tfprotov6.Schema) tftypes.Type {
	bm := map[string]tftypes.Type{}

	for _, att := range schema.Block.Attributes {
		bm[att.Name] = att.Type
	}

	for _, b := range schema.Block.BlockTypes {
		a := map[string]tftypes.Type{}
		for _, att := range b.Block.Attributes {
			a[att.Name] = att.Type
		}
		bm[b.TypeName] = tftypes.List{
			ElementType: tftypes.Object{AttributeTypes: a},
		}

		// FIXME we can make this function recursive to handle
		// n levels of nested blocks
		for _, bb := range b.Block.BlockTypes {
			aa := map[string]tftypes.Type{}
			for _, att := range bb.Block.Attributes {
				aa[att.Name] = att.Type
			}
			a[bb.TypeName] = tftypes.List{
				ElementType: tftypes.Object{AttributeTypes: aa},
			}
		}
	}

	return tftypes.Object{AttributeTypes: bm}
}

// GetResourceType returns the tftypes.Type of a resource of type 'name'
func GetResourceType(name string) (tftypes.Type, error) {
	sch := GetProviderResourceSchema()
	rsch, ok := sch[name]

	if !ok {
		return tftypes.DynamicPseudoType, fmt.Errorf("unknown resource %s - cannot find schema", name)
	}

	return GetObjectTypeFromSchema(rsch), nil
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
