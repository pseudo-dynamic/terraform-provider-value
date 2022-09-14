package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type UnknownProposerResource struct {
}

type UnknownProposerResourceWithTraits interface {
	resource.ResourceWithMetadata
	resource.ResourceWithGetSchema
	resource.ResourceWithModifyPlan
}

func NewUnknownProposerResource() UnknownProposerResourceWithTraits {
	return &UnknownProposerResource{}
}

func (r UnknownProposerResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Version: 0,
		Description: "This resource is very obscure and misbehaving and you really should only use " +
			"it for `value_is_known.proposed_unknown` or `value_is_fully_known.proposed_unknown`.",
		Attributes: map[string]tfsdk.Attribute{
			"value": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "This value will **always** be unknown during the plan phase but always true after apply phase.",
			},
		},
	}, nil
}

type UnknownProposerState struct {
	Value types.Bool `tfsdk:"value"`
}

func (r *UnknownProposerResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_unknown_proposer"
}

// ModifyPlan implements UnknownProposerResourceWithTraits
func (r *UnknownProposerResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Config.Raw.IsNull() {
		// Ignore due to resource deletion
		return
	}

	var plan UnknownProposerState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Value = types.Bool{Unknown: true}
	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create a new resource
func (r UnknownProposerResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan UnknownProposerState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Value = types.Bool{Value: true}
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Read resource information
func (r UnknownProposerResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state UnknownProposerState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

// Update resource
func (r UnknownProposerResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan
	var plan UnknownProposerState
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	plan.Value = types.Bool{Value: true}
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete resource
func (r UnknownProposerResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state UnknownProposerState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}
