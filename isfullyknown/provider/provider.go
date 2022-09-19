package provider

import (
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"

	"github.com/pseudo-dynamic/terraform-provider-value/isknown/common"
	isknown "github.com/pseudo-dynamic/terraform-provider-value/isknown/provider"
)

func Provider() func() tfprotov6.ProviderServer {
	return isknown.ProviderConstructor(isknown.ProviderParameters{
		CheckFullyKnown: true,
	}, common.ProviderResourceSchemaParameters{
		ResourceName: "value_is_fully_known",
		ResourceDescription: "Allows you to have a access to `result` during plan phase that " +
			"states whether `value` or any nested attribute is marked as \"(known after apply)\" or not.",
		ValueDescription: "The `value` and if existing, nested attributes, are tested against \"(known after apply)\"",
		ResultDescription: "States whether `value` or any nested attribute is marked as \"(known after apply)\" or not. If `value` is an aggregate " +
			"type, not only the top level of the aggregate type is checked; elements and attributes " +
			"are checked too.",
	})
}
