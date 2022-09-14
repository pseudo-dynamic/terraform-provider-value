package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
)

type ProviderWithTraits interface {
	provider.ProviderWithMetadata
	provider.ProviderWithResources
}

// Provider schema struct
type ProviderConfig struct {
}

type Provider struct {
}

func NewProvider() ProviderWithTraits {
	return &Provider{}
}

// Metadata implements ProviderWithTraits
func (p *Provider) Metadata(ctx context.Context, req provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "value"
}

// GetSchema
func (p *Provider) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Attributes: map[string]tfsdk.Attribute{},
	}, nil
}

func (p *Provider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	// Retrieve provider data from configuration
	var config ProviderConfig
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// GetResources - Defines provider resources
func (p *Provider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		func() resource.Resource {
			return NewReplacedWhenResource()
		},
	}
}

// // GetDataSources - Defines provider data sources
// func (p *providerType) GetDataSources(_ context.Context) (map[string]datasource.DataSource, diag.Diagnostics) {
// 	return map[string]datasource.DataSource{}, nil
// }
