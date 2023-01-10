package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"
	"terraform-provider-haproxy-pf/haproxy/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &frontendResource{}
	_ resource.ResourceWithConfigure   = &frontendResource{}
	_ resource.ResourceWithImportState = &frontendResource{}
)

// NewFrontendResource is a helper function to simplify the provider implementation.
func NewFrontendResource() resource.Resource {
	return &frontendResource{}
}

// frontendResource is the resource implementation.
type frontendResource struct {
	client *middleware.Client
}

// frontendsModel maps frontends schema data.
type frontendResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Mode               types.String `tfsdk:"mode"`
	Maxconn            types.Int64  `tfsdk:"maxconn"`
	DefaultBackend     types.String `tfsdk:"default_backend"`
	HTTPConnectionMode types.String `tfsdk:"http_connection_mode"`
}

// Metadata returns the resource type name.
func (r *frontendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_frontend"
}

// Schema defines the schema for the resource.
func (r *frontendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"mode": schema.StringAttribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.String{
					middleware.StringDefaultValue(types.StringValue("http")),
				},
			},
			"maxconn": schema.Int64Attribute{
				Optional: true,
				Computed: true,
				PlanModifiers: []planmodifier.Int64{
					middleware.Int64DefaultValue(types.Int64Value(0)),
				},
			},
			"default_backend": schema.StringAttribute{
				Optional: true,
			},
			"http_connection_mode": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Description: "possible values: httpclose,http-server-close,http-keep-alive",
				PlanModifiers: []planmodifier.String{
					middleware.StringDefaultValue(types.StringValue("http-keep-alive")),
				},
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *frontendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*middleware.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *frontendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan frontendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Frontend{
		Name:               plan.Name.ValueString(),
		Mode:               plan.Mode.ValueString(),
		Maxconn:            plan.Maxconn.ValueInt64(),
		DefaultBackend:     plan.DefaultBackend.ValueString(),
		HTTPConnectionMode: plan.HTTPConnectionMode.ValueString(),
	}

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}
	// Create new frontend
	response, err := r.client.CreateFrontend(transaction.Id, payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating frontend",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}
	// commit transaction
	_, err = r.client.CommitTransaction(transaction.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error committing transaction",
			"Could not commit transaction, unexpected error: "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId("root", response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Mode = types.StringValue(response.Mode)
	plan.Maxconn = types.Int64Value(response.Maxconn)
	plan.HTTPConnectionMode = types.StringValue(response.HTTPConnectionMode)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *frontendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state frontendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, frontendName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Get refreshed frontend
	response, err := r.client.GetFrontend(frontendName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Frontend",
			"Could not read Haproxy Frontend ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	id := middleware.CreateResourceId("root", response.Name)
	state.ID = types.StringValue(id)
	state.Name = types.StringValue(response.Name)
	state.Mode = types.StringValue(response.Mode)
	state.Maxconn = types.Int64Value(response.Maxconn)
	state.HTTPConnectionMode = types.StringValue(response.HTTPConnectionMode)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *frontendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from state. extract ID
	var state frontendResourceModel
	req.State.Get(ctx, &state)

	// Retrieve values from plan
	var plan frontendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Frontend{
		Name:               plan.Name.ValueString(),
		Mode:               plan.Mode.ValueString(),
		Maxconn:            plan.Maxconn.ValueInt64(),
		DefaultBackend:     plan.DefaultBackend.ValueString(),
		HTTPConnectionMode: plan.HTTPConnectionMode.ValueString(),
	}
	_, frontendName, err := middleware.ResourceParseId(ctx, state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting resource ID",
			"Could not update frontend, unexpected error: "+err.Error(),
		)
		return
	}

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing frontend
	_, err = r.client.UpdateFrontend(transaction.Id, frontendName, payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Haproxy Frontend",
			"Could not update frontend, unexpected error: "+err.Error(),
		)
		return
	}

	// commit transaction
	_, err = r.client.CommitTransaction(transaction.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error committing transaction",
			"Could not commit transaction, unexpected error: "+err.Error(),
		)
		return
	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	response, err := r.client.GetFrontend(frontendName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Frontend",
			"Could not read Haproxy Frontend ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId("root", response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Mode = types.StringValue(response.Mode)
	plan.Maxconn = types.Int64Value(response.Maxconn)
	plan.HTTPConnectionMode = types.StringValue(response.HTTPConnectionMode)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *frontendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state frontendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, frontendName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create frontend, unexpected error: "+err.Error(),
		)
		return
	}

	// Delete existing frontend
	err = r.client.DeleteFrontend(transaction.Id, frontendName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Haproxy Frontend",
			"Could not delete frontend, unexpected error: "+err.Error(),
		)
		return
	}

	// commit transaction
	_, err = r.client.CommitTransaction(transaction.Id)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error committing transaction",
			"Could not commit transaction, unexpected error: "+err.Error(),
		)
		return
	}

}

func (r *frontendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	_, frontendName, err := middleware.ResourceParseId(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Cannot parse import ID, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), frontendName)...)
}
