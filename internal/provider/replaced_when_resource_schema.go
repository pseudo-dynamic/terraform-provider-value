package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

func (r ReplacedWhenResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `Enables the scenario to only change a value when a condition is met. 
The value attribute can for example be used as target for replace_triggered_by. To detect
resource creation and resource deletion as change you can try the following approach:

	resource "value_replaced_when" "true" {
		count     = 1
		condition = true
	}

	resource "value_stash" "replacement_trigger" {
		// This workaround detects not only resource creation
		// and every change of value but also the deletion of
		// resource.
		value = try(value_replaced_when.true[0].value, null)
	}

	resource "value_stash" "replaced" {
		lifecycle {
			// Replace me whenever value_replaced_when.true[0].value 
			// changes or value_replaced_when[0] gets deleted.
			replace_triggered_by = [
				value_stash.replacement_trigger.value
			]
		}
	}`,
		Attributes: map[string]tfsdk.Attribute{
			"condition": {
				Type:                types.BoolType,
				Required:            true,
				Description:         "If already true or getting true then value will be replaced by a random value",
				MarkdownDescription: "If already `true` or getting `true` then `value` will be replaced by a random value",
			},
			"value": {
				Type:     types.StringType,
				Computed: true,
				Description: `If the very first condition is false or remains false, then the value remains unchanged. 
The condition change from true to false does not trigger a replacement of those who use the value as
target for replace_triggered_by.

If the condition is true or remains true, then the value will be always updated in-place with a random
value. It will always trigger a replacement of those who use the value as target for replace_triggered_by.

There is a special case that a replacement of those who use the value as target for replace_triggered_by
occurs when the condition is unknown and uncomputed. This always happens whenever the condition becomes
unknown and uncomputed again.`,
			},
		},
	}, nil
}

// type ValuePlanModifier struct{}

// func (ValuePlanModifier) Description(ctx context.Context) string {
// 	return ""
// }

// func (ValuePlanModifier) MarkdownDescription(ctx context.Context) string {
// 	return ""
// }

// func (ValuePlanModifier) Modify(ctx context.Context, req tfsdk.ModifyAttributePlanRequest, resp *tfsdk.ModifyAttributePlanResponse) {
// }
