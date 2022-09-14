package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// getResourceType returns the tftypes.Type of a resource of type 'name'
func getResourceType(name string) tftypes.Type {
	sch := getProviderResourceSchema(name)
	rsch, ok := sch[name]

	if !ok {
		panic(fmt.Errorf("unknown resource %s - cannot find schema", name))
	}

	return getObjectTypeFromSchema(rsch)
}

// getObjectTypeFromSchema returns a tftypes.Type that can wholy represent the schema input
func getObjectTypeFromSchema(schema *tfprotov6.Schema) tftypes.Type {
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

// getDocumentedProviderResourceSchema contains the definitions of all supported resources with documentations
func getDocumentedProviderResourceSchema(params common.ProviderResourceSchemaParameters) map[string]*tfprotov6.Schema {
	return common.GetProviderResourceSchema(params)
}

// getProviderResourceSchema contains the definitions of all supported resources
func getProviderResourceSchema(resourceName string) map[string]*tfprotov6.Schema {
	return getDocumentedProviderResourceSchema(common.ProviderResourceSchemaParameters{
		ResourceName: resourceName,
	})
}
