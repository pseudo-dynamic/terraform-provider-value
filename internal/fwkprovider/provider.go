package fwkprovider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	fwkprovider "github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/pseudo-dynamic/terraform-provider-value/internal/fwkproviderconfig"
)

const providerName string = "value"

type provider struct {
}

type providerWithTraits interface {
	fwkprovider.ProviderWithMetadata
	fwkprovider.ProviderWithResources
	fwkprovider.ProviderWithDataSources
	fwkprovider.ProviderWithMetaSchema
}

// Provider schema struct
type providerConfig struct {
	GuidSeedAddition types.String `tfsdk:"guid_seed_addition"`
}

type providerData struct {
	GuidSeedAddition *string
}

func NewProvider() providerWithTraits {
	return &provider{}
}

// Metadata implements ProviderWithTraits
func (p *provider) Metadata(ctx context.Context, req fwkprovider.MetadataRequest, resp *fwkprovider.MetadataResponse) {
	resp.TypeName = "value"
}

// GetSchema
func (p *provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return *fwkproviderconfig.GetProviderConfigSchema(), nil
}

// GetMetaSchema implements providerWithTraits
func (*provider) GetMetaSchema(context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return *fwkproviderconfig.GetProviderMetaSchema(), nil
}

func (p *provider) Configure(ctx context.Context, req fwkprovider.ConfigureRequest, resp *fwkprovider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config providerConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}

	resp.ResourceData = &providerData{
		GuidSeedAddition: &config.GuidSeedAddition.Value,
	}
}

// GetResources - Defines provider resources
func (p *provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return NewReplacedWhenResource()
		},
		func() resource.Resource {
			return NewUnknownProposerResource()
		},
		func() resource.Resource {
			return NewPathExistsResource()
		},
	}
}

// DataSources implements ProviderWithTraits
func (*provider) DataSources(context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		func() datasource.DataSource {
			return NewTempDirDataSource()
		},
	}
}
