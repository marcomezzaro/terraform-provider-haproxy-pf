package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"
	"terraform-provider-haproxy-pf/haproxy/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/avast/retry-go/v4"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &bindResource{}
	_ resource.ResourceWithConfigure   = &bindResource{}
	_ resource.ResourceWithImportState = &bindResource{}
)

// NewBindResource is a helper function to simplify the provider implementation.
func NewBindResource() resource.Resource {
	return &bindResource{}
}

// bindResource is the resource implementation.
type bindResource struct {
	client *middleware.Client
}

// bindsModel maps binds schema data.
type bindResourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	Address    types.String `tfsdk:"address"`
	Port       types.Int64  `tfsdk:"port"`
	ParentName types.String `tfsdk:"parent_name"`
}

// Metadata returns the resource type name.
func (r *bindResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_bind"
}

// Schema defines the schema for the resource.
func (r *bindResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
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
func (r *bindResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*middleware.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *bindResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan bindResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Bind{
		Name:    plan.Name.ValueString(),
		Address: plan.Address.ValueString(),
		Port:    plan.Port.ValueInt64(),
	}

	var response *models.Bind
	retry_err := retry.Do(
		func() error {
			// Open transaction
			configuration, err := r.client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := r.client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}
			// Create new bind
			create_response, err := r.client.CreateBind(transaction.Id, payload, plan.ParentName.ValueString())
			if err != nil {
				return nil
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			response = create_response
			return nil
		})

	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error creating bind",
			"Could not create bind, unexpected error: "+retry_err.Error(),
		)
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId(plan.ParentName.ValueString(), response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Address = types.StringValue(response.Address)
	plan.Port = types.Int64Value(response.Port)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *bindResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state bindResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, bindName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Get refreshed bind
	response, err := r.client.GetBind(bindName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Bind",
			"Could not read Haproxy Bind ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	id := middleware.CreateResourceId(parentName, response.Name)
	state.ID = types.StringValue(id)
	state.Name = types.StringValue(response.Name)
	state.Address = types.StringValue(response.Address)
	state.Port = types.Int64Value(response.Port)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *bindResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from state. extract ID
	var state bindResourceModel
	req.State.Get(ctx, &state)

	// Retrieve values from plan
	var plan bindResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.Bind{
		Name:    plan.Name.ValueString(),
		Address: plan.Address.ValueString(),
		Port:    plan.Port.ValueInt64(),
	}
	parentName, bindName, err := middleware.ResourceParseId(ctx, state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting resource ID",
			"Could not update bind, unexpected error: "+err.Error(),
		)
		return
	}

	retry_err := retry.Do(
		func() error {
			// Open transaction
			configuration, err := r.client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := r.client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}
			// Update existing bind
			_, err = r.client.UpdateBind(transaction.Id, payload, parentName)
			if err != nil {
				return err
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		})

	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error updating bind",
			"Could not update bind, unexpected error: "+retry_err.Error(),
		)
		return
	}

	response, err := r.client.GetBind(bindName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Bind",
			"Could not read Haproxy Bind ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId(parentName, response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Address = types.StringValue(response.Address)
	plan.Port = types.Int64Value(response.Port)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *bindResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state bindResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, bindName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	retry_err := retry.Do(
		func() error {
			// Open transaction
			configuration, err := r.client.GetConfiguration()
			if err != nil {
				return err
			}
			transaction, err := r.client.CreateTransaction(configuration.Version)
			if err != nil {
				return err
			}
			// Delete existing bind
			err = r.client.DeleteBind(transaction.Id, bindName, parentName)
			if err != nil {
				return err
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		})
	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error deleting bind",
			"Could not delete bind, unexpected error: "+retry_err.Error(),
		)
		return
	}

}

func (r *bindResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	parentName, bindName, err := middleware.ResourceParseId(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Cannot parse import ID, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), bindName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("parent_name"), parentName)...)
}
