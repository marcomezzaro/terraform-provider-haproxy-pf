package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"
	"terraform-provider-haproxy-pf/haproxy/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &serverResource{}
	_ resource.ResourceWithConfigure   = &serverResource{}
	_ resource.ResourceWithImportState = &serverResource{}
)

// NewServerResource is a helper function to simplify the provider implementation.
func NewServerResource() resource.Resource {
	return &serverResource{}
}

// serverResource is the resource implementation.
type serverResource struct {
	client *middleware.Client
}

// serversModel maps servers schema data.
type serverResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Address    types.String `tfsdk:"address"`
	Check      types.String `tfsdk:"check"`
	Port       types.Int64  `tfsdk:"port"`
	ParentName types.String `tfsdk:"parent_name"`
}

// Metadata returns the resource type name.
func (r *serverResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server"
}

// Schema defines the schema for the resource.
func (r *serverResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"address": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"check": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"port": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"parent_name": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *serverResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*middleware.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *serverResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan serverResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Server{
		Name:    plan.Name.ValueString(),
		Address: plan.Address.ValueString(),
		Port:    plan.Port.ValueInt64(),
		Check:   plan.Check.ValueString(),
	}

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}
	// Create new server
	response, err := r.client.CreateServer(transaction.Id, payload, plan.ParentName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating server",
			"Could not create server, unexpected error: "+err.Error(),
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
	id := middleware.CreateResourceId(plan.ParentName.ValueString(), response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Address = types.StringValue(response.Address)
	plan.Port = types.Int64Value(response.Port)
	plan.Check = types.StringValue(response.Check)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serverResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state serverResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, serverName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Get refreshed server
	response, err := r.client.GetServer(serverName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Server",
			"Could not read Haproxy Server ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	id := middleware.CreateResourceId(parentName, response.Name)
	state.ID = types.StringValue(id)
	state.Name = types.StringValue(response.Name)
	state.Address = types.StringValue(response.Address)
	state.Port = types.Int64Value(response.Port)
	state.Check = types.StringValue(response.Check)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from state. extract ID
	var state serverResourceModel
	req.State.Get(ctx, &state)

	// Retrieve values from plan
	var plan serverResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Server{
		Name:    plan.Name.ValueString(),
		Address: plan.Address.ValueString(),
		Port:    plan.Port.ValueInt64(),
		Check:   plan.Check.ValueString(),
	}
	parentName, serverName, err := middleware.ResourceParseId(ctx, state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting resource ID",
			"Could not update server, unexpected error: "+err.Error(),
		)
		return
	}

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}

	// Update existing server
	_, err = r.client.UpdateServer(transaction.Id, payload, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Haproxy Server",
			"Could not update server, unexpected error: "+err.Error(),
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

	response, err := r.client.GetServer(serverName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Server",
			"Could not read Haproxy Server ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId(parentName, response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Address = types.StringValue(response.Address)
	plan.Port = types.Int64Value(response.Port)
	plan.Check = types.StringValue(response.Check)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serverResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state serverResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, serverName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Open transaction
	configuration, err := r.client.GetConfiguration()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting configuration",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}
	transaction, err := r.client.CreateTransaction(configuration.Version)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating haproxy dataplane transaction",
			"Could not create server, unexpected error: "+err.Error(),
		)
		return
	}

	// Delete existing server
	err = r.client.DeleteServer(transaction.Id, serverName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Haproxy Server",
			"Could not delete server, unexpected error: "+err.Error(),
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

func (r *serverResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	parentName, serverName, err := middleware.ResourceParseId(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Cannot parse import ID, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), serverName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("parent_name"), parentName)...)
}
