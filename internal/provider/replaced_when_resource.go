package provider

import (
	"context"
	"hash/fnv"
	"math/rand"

	"github.com/google/uuid"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type ReplacedWhenResourceWithTraits interface {
	resource.ResourceWithMetadata
	resource.ResourceWithGetSchema
	resource.ResourceWithModifyPlan
}

type ReplacedWhenState struct {
	When  types.Bool   `tfsdk:"condition"`
	Value types.String `tfsdk:"value"`
}

type ReplacedWhenResource struct {
}

func NewReplacedWhenResource() ReplacedWhenResourceWithTraits {
	return &ReplacedWhenResource{}
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

	if !suppliedCondition.Unknown {
		if currentValue.IsUnknown() || suppliedCondition.Value {
			// First creation of attribute / state value won't never be null again
			// OR supplied condition (config) is known (the latter condition is also
			// true if the unknown condition has been computed once and didn't change
			// and is therefore not unknown anymore)

			// if currentValue.IsUnknown() {
			// 	newValue = types.String{Null: true}
			// } else {
			// Because ModifyPlan gets called twice but without any shared context
			// we need to create deterministic incremental UUIDs that won't change
			// in apply and its following plan phase.
			// I chose this approach because I want this resource to be as minimalistic
			// as possible. The alternative would be to work with two pre-chosen values
			// and swap between them back and forth.
			hash := fnv.New64a()
			hash.Write([]byte(currentValue.Value))
			hashUnsignedSum := hash.Sum64()
			hashSignedSum := int64(hashUnsignedSum)
			rnd := rand.New(rand.NewSource(hashSignedSum))
			deterministicUuid, _ := uuid.NewRandomFromReader(rnd)
			nextUuid := deterministicUuid.String()
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

	resp.State = req.State
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
