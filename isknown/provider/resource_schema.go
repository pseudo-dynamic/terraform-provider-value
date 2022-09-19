package provider

import (
	"fmt"

	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/schema"
	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
)

// getResourceType returns the tftypes.Type of a resource of type 'name'
func getResourceType(name string) tftypes.Type {
	sch := getProviderResourceSchema(name)
	rsch, ok := sch[name]

	if !ok {
		panic(fmt.Errorf("unknown resource %s - cannot find schema", name))
	}

	return schema.GetObjectTypeFromSchema(rsch)
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
