package provider

import (
	"context"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type TempDirDataSource struct {
}

type TempDirDataSourceWithTraits interface {
	datasource.DataSource
}

func NewTempDirDataSource() TempDirDataSourceWithTraits {
	return &TempDirDataSource{}
}

func (r TempDirDataSource) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: "Simply returns the OS-dependent temporary directory (e.g. /tmp).",
		Attributes: map[string]tfsdk.Attribute{
			"path": {
				Type:        types.StringType,
				Computed:    true,
				Description: "The OS-dependent temporary directory.",
			},
		},
	}, nil
}

func (r *TempDirDataSource) Metadata(ctx context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_temp_dir"
}

type TempDirState struct {
	Path types.String `tfsdk:"path"`
}

// Read resource information
func (r TempDirDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Get current state
	var state TempDirState
	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	state.Path = types.String{Value: os.TempDir()}
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}
