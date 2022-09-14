package provider

import (
	"context"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/uuid"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReplacedWhenResource struct {
}

type ReplacedWhenResourceWithTraits interface {
	resource.ResourceWithMetadata
	resource.ResourceWithGetSchema
	resource.ResourceWithModifyPlan
}

func NewReplacedWhenResource() ReplacedWhenResourceWithTraits {
	return &ReplacedWhenResource{}
}

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
				Description: `If the very first condition is false, then the value will be once initialized by a random value.

If the condition is false or remains false, then the value remains unchanged. 
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

type ReplacedWhenState struct {
	When  types.Bool   `tfsdk:"condition"`
	Value types.String `tfsdk:"value"`
}

func (r *ReplacedWhenResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_replaced_when"
}

// ModifyPlan implements ReplacedWhenResourceWithTraits
func (r *ReplacedWhenResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		// Ignore due to resource deletion
		return
	}

	var plannedValue types.String
	diags := req.Plan.GetAttribute(ctx, path.Root("value"), &plannedValue)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// resp.Diagnostics.AddWarning("Value (plan) is known: "+strconv.FormatBool(!plannedValue.IsUnknown()), "")
	// resp.Diagnostics.AddWarning("Value (plan) is null: "+strconv.FormatBool(plannedValue.IsNull()), "")
	// resp.Diagnostics.AddWarning("Value (plan) is: "+plannedValue.Value, "")

	var currentValue types.String
	diags = req.State.GetAttribute(ctx, path.Root("value"), &currentValue)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// if plannedValue.Unknown {
	var suppliedCondition types.Bool
	diags = req.Config.GetAttribute(ctx, path.Root("condition"), &suppliedCondition)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	// resp.Diagnostics.AddWarning("Condition (config) is known: "+strconv.FormatBool(!suppliedCondition.IsUnknown()), "")
	// resp.Diagnostics.AddWarning("Condition (config) is null: "+strconv.FormatBool(suppliedCondition.IsNull()), "")
	// resp.Diagnostics.AddWarning("Condition (config) is: "+strconv.FormatBool(suppliedCondition.Value), "")

	var suppliedValue types.String
	diags = req.Config.GetAttribute(ctx, path.Root("value"), &suppliedValue)

	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	var newValue types.String
	var currentValueNotOnceSet = currentValue.IsUnknown() || currentValue.IsNull()

	if currentValueNotOnceSet || !suppliedCondition.Unknown {
		// Empty value (state) / state value won't never be unknown or null again
		// OR supplied condition (config) is known (the latter condition is also
		// true if the unknown condition has been computed once and didn't change
		// and is therefore not unknown anymore)
		if currentValueNotOnceSet || suppliedCondition.Value {
			// if currentValue.IsUnknown() {
			// 	newValue = types.String{Null: true}
			// } else {
			// Because ModifyPlan gets called twice but without any shared context
			// we need to create deterministic incremental UUIDs that won't change
			// in apply and its following plan phase.
			// I chose this approach because I want this resource to be as minimalistic
			// as possible. The alternative would be to work with two pre-chosen values
			// and swap between them back and forth.
			nextUuid := uuid.DeterministicUuidFromString(currentValue.Value).String()
			newValue = types.String{Value: nextUuid}
			// }
		} else {
			diags := req.State.GetAttribute(ctx, path.Root("value"), &newValue)
			resp.Diagnostics.Append(diags...)

			if resp.Diagnostics.HasError() {
				return
			}
		}
	} else {
		newValue = types.String{Unknown: true}
	}

	newPlan := req.Plan
	newPlan.SetAttribute(ctx, path.Root("value"), &newValue)
	resp.Plan = newPlan
}

// Create a new resource
func (r ReplacedWhenResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan ReplacedWhenState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Read resource information
func (r ReplacedWhenResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state ReplacedWhenState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// resp.State = req.State
}

// Update resource
func (r ReplacedWhenResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan
	var plan ReplacedWhenState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Set new state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete resource
func (r ReplacedWhenResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state ReplacedWhenState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}
