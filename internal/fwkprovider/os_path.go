package fwkprovider

import (
	"context"
	"os"

	"github.com/pseudo-dynamic/terraform-provider-value/internal/fwkproviderconfig"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/goproviderconfig"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/guid"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

const osPathResourceSuffix string = "_os_path"
const osPathResourceName string = providerName + osPathResourceSuffix

type osPathResource struct {
	ProviderGuidSeedAddition *string
}

type osPathResourceWithTraits interface {
	resource.ResourceWithMetadata
	resource.ResourceWithGetSchema
	resource.ResourceWithModifyPlan
	resource.ResourceWithConfigure
}

func NewOSPathResource() osPathResourceWithTraits {
	return &osPathResource{}
}

// Configure implements osPathResourceWithTraits
func (r *osPathResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		// For whatever reason Configure gets called with nil ProviderData.
		return
	}

	provderData := req.ProviderData.(*providerData)
	r.ProviderGuidSeedAddition = provderData.GuidSeedAddition
}

func (r osPathResource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Checks if an OS path exists and caches its computation at plan-time and won't change after " +
			"apply-time even the path may have been removed." + "\n" + goproviderconfig.GetProviderMetaGuidSeedAdditionAttributeDescription(),
		Attributes: map[string]tfsdk.Attribute{
			"path": {
				Type:        types.StringType,
				Required:    true,
				Description: "A path to a file or directory.",
			},
			"proposed_unknown": {
				Type:        types.BoolType,
				Required:    true,
				Description: goproviderconfig.GetProposedUnknownAttributeDescription(),
			},
			"guid_seed": {
				Type:        types.StringType,
				Required:    true,
				Description: goproviderconfig.GetGuidSeedAttributeDescription(osPathResourceName),
			},
			"exists": {
				Type:        types.BoolType,
				Computed:    true,
				Description: "The computation whether the path exists or not.",
			},
		},
	}, nil
}

type osPathState struct {
	Path            types.String `tfsdk:"path"`
	GuidSeed        types.String `tfsdk:"guid_seed"`
	ProposedUnknown types.Bool   `tfsdk:"proposed_unknown"`
	Exists          types.Bool   `tfsdk:"exists"`
}

func (r *osPathResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + osPathResourceSuffix
}

// ModifyPlan implements OSPathResourceWithTraits
func (r *osPathResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if r.ProviderGuidSeedAddition == nil {
		resp.Diagnostics.AddError("Bad provider guid seed", "Provider guid seed is null but was expected to be empty")
		return
	}

	// Get current config
	var config osPathState
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
	suppliedGuidSeed := config.GuidSeed.Value
	isPlanPhase := config.ProposedUnknown.IsUnknown()

	if !fwkproviderconfig.ValidatePlanKnownString(config.GuidSeed, "guid_seed", &resp.Diagnostics) {
		return
	}

	var providerMetaSeedAddition string
	var isSuccessful bool
	if providerMetaSeedAddition, _, isSuccessful = goproviderconfig.TryUnmarshalValueThenExtractGuidSeedAddition(&req.ProviderMeta.Raw); !isSuccessful {
		resp.Diagnostics.AddError("Extraction failed", "Could not extract provider meta guid seed addition")
		return
	}

	composedGuidSeed := guid.ComposeGuidSeed(r.ProviderGuidSeedAddition,
		&providerMetaSeedAddition,
		osPathResourceName,
		"exists",
		&suppliedGuidSeed)

	checkPathExistence := func() types.Bool {
		if config.Path.Unknown {
			return types.Bool{Unknown: true}
		}

		_, err := os.Stat(config.Path.Value)
		osPath := err == nil
		return types.Bool{Value: osPath}
	}

	cachedExists, err := guid.GetPlanCachedBoolean(
		isPlanPhase,
		composedGuidSeed,
		osPathResourceName,
		checkPathExistence)

	if err != nil {
		resp.Diagnostics.AddError("Plan cache mechanism failed for exists attribute", err.Error())
		return
	}

	config.Exists = cachedExists
	resp.Plan.Set(ctx, &config)
}

// Create a new resource
func (r osPathResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan osPathState
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
func (r osPathResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state osPathState
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
}

// Update resource
func (r osPathResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan
	var plan osPathState
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
func (r osPathResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state osPathState
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	// Remove resource from state
	resp.State.RemoveResource(ctx)
}
