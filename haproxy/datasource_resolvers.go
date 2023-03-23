package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &resolversDataSource{}
	_ datasource.DataSourceWithConfigure = &resolversDataSource{}
)

func NewResolversDataSource() datasource.DataSource {
	return &resolversDataSource{}
}

type resolversDataSource struct {
	client *middleware.Client
}

func (d *resolversDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_resolvers"
}

// Configure adds the provider configured client to the data source.
func (d *resolversDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*middleware.Client)
}

func (d *resolversDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"resolvers": schema.ListNestedAttribute{
				Computed: true,
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"id": schema.StringAttribute{
							Computed: true,
						},
						"name": schema.StringAttribute{
							Computed: true,
						},
					},
				},
			},
		},
	}
}

// resolversDataSourceModel maps the data source schema data.
type resolversDataSourceModel struct {
	Resolvers []resolversModel `tfsdk:"resolvers"`
}

// resolversModel maps resolvers schema data.
type resolversModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
}

func (d *resolversDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {

	var state resolversDataSourceModel

	resolvers, err := d.client.GetResolvers()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Haproxy Resolvers",
			err.Error(),
		)
		return
	}

	// Map response body to model
	for _, resolver := range resolvers.Data {
		resolverState := resolversModel{
			ID:      types.StringValue(resolver.Name),
			Name:    types.StringValue(resolver.Name),
		}

		state.Resolvers = append(state.Resolvers, resolverState)
	}

	// Set state
	diags := resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}
