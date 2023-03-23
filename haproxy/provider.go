package haproxy

import (
	"context"
	"os"

	"terraform-provider-haproxy-pf/haproxy/middleware"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces
var (
	_ provider.Provider = &haproxyProvider{}
)

// New is a helper function to simplify provider server and testing implementation.
func New() provider.Provider {
	return &haproxyProvider{}
}

// haproxyProvider is the provider implementation.
type haproxyProvider struct{}

// haproxyProviderModel maps provider schema data to a Go type.
type haproxyProviderModel struct {
	Host     types.String `tfsdk:"host"`
	Username types.String `tfsdk:"username"`
	Password types.String `tfsdk:"password"`
	Insecure types.Bool   `tfsdk:"insecure"`
}

// Metadata returns the provider type name.
func (p *haproxyProvider) Metadata(_ context.Context, _ provider.MetadataRequest, resp *provider.MetadataResponse) {
	resp.TypeName = "haproxy-pf"
}

// Schema defines the provider-level schema for configuration data.
func (p *haproxyProvider) Schema(_ context.Context, _ provider.SchemaRequest, resp *provider.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"host": schema.StringAttribute{
				Optional: true,
			},
			"username": schema.StringAttribute{
				Optional: true,
			},
			"password": schema.StringAttribute{
				Optional:  true,
				Sensitive: true,
			},
			"insecure": schema.BoolAttribute{
				Optional:  true,
				Sensitive: false,
			},
		},
	}

}

// Configure prepares a haproxy API client for data sources and resources.
func (p *haproxyProvider) Configure(ctx context.Context, req provider.ConfigureRequest, resp *provider.ConfigureResponse) {
	tflog.Info(ctx, "Configuring Haproxy client")

	// Retrieve provider data from configuration
	var config haproxyProviderModel
	diags := req.Config.Get(ctx, &config)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if config.Host.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"), "Unknown haproxy API Host", "")
	}

	if config.Username.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"), "Unknown haproxy API Username", "")
	}

	if config.Password.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"), "Unknown haproxy API Password", "")
	}

	if config.Insecure.IsUnknown() {
		resp.Diagnostics.AddAttributeError(
			path.Root("insecure"), "Unknown haproxy API Insecure flag", "")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// Default values to environment variables, but override
	// with Terraform configuration value if set.

	host := os.Getenv("HAPROXY_HOST")
	username := os.Getenv("HAPROXY_USERNAME")
	password := os.Getenv("HAPROXY_PASSWORD")
	insecure := false

	if !config.Host.IsNull() {
		host = config.Host.ValueString()
	}

	if !config.Username.IsNull() {
		username = config.Username.ValueString()
	}

	if !config.Password.IsNull() {
		password = config.Password.ValueString()
	}

	if !config.Insecure.IsNull() {
		insecure = config.Insecure.ValueBool()
	}

	// If any of the expected configurations are missing, return
	// errors with provider-specific guidance.

	if host == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("host"), "Unknown haproxy API Host", "")
	}

	if username == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("username"), "Unknown haproxy API Username", "")
	}

	if password == "" {
		resp.Diagnostics.AddAttributeError(
			path.Root("password"), "Unknown haproxy API Password", "")

	}

	if resp.Diagnostics.HasError() {
		return
	}

	ctx = tflog.SetField(ctx, "haproxy_host", host)
	ctx = tflog.SetField(ctx, "haproxy_username", username)
	ctx = tflog.SetField(ctx, "haproxy_password", password)
	ctx = tflog.SetField(ctx, "haproxy_insecure", insecure)
	ctx = tflog.MaskFieldValuesWithFieldKeys(ctx, "haproxy_password")

	tflog.Info(ctx, "Creating Haproxy client")

	// Create a new Haproxy client using the configuration values
	client := middleware.NewClient(username, password, host, insecure)

	// Make the Haproxy client available during DataSource and Resource
	// type Configure methods.
	resp.DataSourceData = client
	resp.ResourceData = client

	tflog.Info(ctx, "Configured Haproxy client", map[string]any{"success": true})

}

// DataSources defines the data sources implemented in the provider.
func (p *haproxyProvider) DataSources(_ context.Context) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		NewBackendsDataSource,
		NewFrontendsDataSource,
		NewResolversDataSource,
	}
}

// Resources defines the resources implemented in the provider.
func (p *haproxyProvider) Resources(_ context.Context) []func() resource.Resource {
	return []func() resource.Resource{
		NewBackendResource,
		NewFrontendResource,
		NewBindResource,
		NewServerResource,
		NewServerTemplateResource,
	}
}
