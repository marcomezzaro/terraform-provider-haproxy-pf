package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &backendsDataSource{}
	_ datasource.DataSourceWithConfigure = &backendsDataSource{}
)

func NewBackendsDataSource() datasource.DataSource {
	return &backendsDataSource{}
}

type backendsDataSource struct {
	client *middleware.Client
}

func (d *backendsDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backends"
}

// Configure adds the provider configured client to the data source.
func (d *backendsDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*middleware.Client)
}

func (d *backendsDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"backends": schema.ListNestedAttribute{
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
							Computed: true,
						},
						"balance": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// backendsDataSourceModel maps the data source schema data.
type backendsDataSourceModel struct {
	Backends []backendsModel `tfsdk:"backends"`
}

// backendsModel maps backends schema data.
type backendsModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Mode    types.String `tfsdk:"mode"`
	Balance types.String `tfsdk:"balance"`
}

func (d *backendsDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state backendsDataSourceModel

	backends, err := d.client.GetBackends()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Haproxy Backends",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, backend := range backends.Data {
		backendState := backendsModel{
			ID:      types.StringValue(backend.Name),
			Name:    types.StringValue(backend.Name),
			Mode:    types.StringValue(backend.Mode),
			Balance: types.StringValue(backend.Balance.Algorithm),
		}

		state.Backends = append(state.Backends, backendState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
