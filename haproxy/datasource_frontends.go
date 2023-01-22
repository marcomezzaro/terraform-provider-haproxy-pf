package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &frontendsDataSource{}
	_ datasource.DataSourceWithConfigure = &frontendsDataSource{}
)

func NewFrontendsDataSource() datasource.DataSource {
	return &frontendsDataSource{}
}

type frontendsDataSource struct {
	client *middleware.Client
}

func (d *frontendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_frontends"
}

// Configure adds the provider configured client to the data source.
func (d *frontendsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*middleware.Client)
}

func (d *frontendsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"frontends": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
						"mode": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"maxconn": schema.Int64Attribute{
							Optional: true,
							Computed: true,
						},
						"default_backend": schema.StringAttribute{
							Optional: true,
							Computed: true,
						},
						"http_connection_mode": schema.StringAttribute{
							Optional:    true,
							Computed:    true,
							Description: "possible values: httpclose,http-server-close,http-keep-alive",
						},
					},
				},
			},
		},
	}
}

// frontendsDataSourceModel maps the data source schema data.
type frontendsDataSourceModel struct {
	Frontends []frontendsModel `tfsdk:"frontends"`
}

// frontendsModel maps frontends schema data.
type frontendsModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Mode               types.String `tfsdk:"mode"`
	Maxconn            types.Int64  `tfsdk:"maxconn"`
	DefaultBackend     types.String `tfsdk:"default_backend"`
	HTTPConnectionMode types.String `tfsdk:"http_connection_mode"`
}

func (d *frontendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state frontendsDataSourceModel

	frontends, err := d.client.GetFrontends()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Haproxy Frontends",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, frontend := range frontends.Data {
		frontendState := frontendsModel{
			ID:                 types.StringValue(frontend.Name),
			Name:               types.StringValue(frontend.Name),
			Mode:               types.StringValue(frontend.Mode),
			Maxconn:            types.Int64Value(frontend.Maxconn),
			HTTPConnectionMode: types.StringValue(frontend.HTTPConnectionMode),
		}

		state.Frontends = append(state.Frontends, frontendState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
