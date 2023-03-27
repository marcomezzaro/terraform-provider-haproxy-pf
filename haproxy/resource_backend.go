package haproxy

import (
	"context"
	"terraform-provider-haproxy-pf/haproxy/middleware"
	"terraform-provider-haproxy-pf/haproxy/models"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/avast/retry-go/v4"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &backendResource{}
	_ resource.ResourceWithConfigure   = &backendResource{}
	_ resource.ResourceWithImportState = &backendResource{}
)

// NewBackendResource is a helper function to simplify the provider implementation.
func NewBackendResource() resource.Resource {
	return &backendResource{}
}

// backendResource is the resource implementation.
type backendResource struct {
	client *middleware.Client
}

// backendsModel maps backends schema data.
type backendResourceModel struct {
	ID      types.String `tfsdk:"id"`
	Name    types.String `tfsdk:"name"`
	Mode    types.String `tfsdk:"mode"`
	Balance types.String `tfsdk:"balance"`
}

// Metadata returns the resource type name.
func (r *backendResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_backend"
}

// Schema defines the schema for the resource.
func (r *backendResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Optional: false,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"mode": schema.StringAttribute{
				Optional: true,
			},
			"balance": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
		},
	}
}

// Configure adds the provider configured client to the resource.
func (r *backendResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*middleware.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *backendResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var balance = models.Balance{
		Algorithm: plan.Balance.ValueString(),
	}
	var payload = models.Backend{
		Name:    plan.Name.ValueString(),
		Mode:    plan.Mode.ValueString(),
		Balance: balance,
	}

	var response *models.Backend
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
			// Create new backend
			create_response, err := r.client.CreateBackend(transaction.Id, payload)
			if err != nil {
				return err
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			response = create_response
			return nil
		},
	)

	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error creating backend",
			"Could not create backend, unexpected error: "+retry_err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId("root", response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Mode = types.StringValue(response.Mode)
	plan.Balance = types.StringValue(response.Balance.Algorithm)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *backendResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, backendName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Get refreshed backend
	response, err := r.client.GetBackend(backendName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Backend",
			"Could not read Haproxy Backend ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	id := middleware.CreateResourceId("root", response.Name)
	state.ID = types.StringValue(id)
	state.Name = types.StringValue(response.Name)
	state.Mode = types.StringValue(response.Mode)
	state.Balance = types.StringValue(response.Balance.Algorithm)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *backendResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from state. extract ID
	var state backendResourceModel
	req.State.Get(ctx, &state)

	// Retrieve values from plan
	var plan backendResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var balance = models.Balance{
		Algorithm: plan.Balance.ValueString(),
	}
	var payload = models.Backend{
		Name:    plan.Name.ValueString(),
		Mode:    plan.Mode.ValueString(),
		Balance: balance,
	}
	_, backendName, err := middleware.ResourceParseId(ctx, state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting resource ID",
			"Could not update backend, unexpected error: "+err.Error(),
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
			// Update existing backend
			_, err = r.client.UpdateBackend(transaction.Id, backendName, payload)
			if err != nil {
				return err
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error updating backend",
			"Could not update backend, unexpected error: "+retry_err.Error(),
		)
		return

	}

	// Fetch updated items from GetOrder as UpdateOrder items are not
	// populated.
	response, err := r.client.GetBackend(backendName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy Backend",
			"Could not read Haproxy Backend ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId("root", response.Name)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Name)
	plan.Mode = types.StringValue(response.Mode)
	plan.Balance = types.StringValue(response.Balance.Algorithm)

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *backendResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state backendResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, backendName, _ := middleware.ResourceParseId(ctx, state.ID.String())

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
			// Delete existing backend
			err = r.client.DeleteBackend(transaction.Id, backendName)
			if err != nil {
				return err
			}
			// commit transaction
			_, err = r.client.CommitTransaction(transaction.Id)
			if err != nil {
				return err
			}
			return nil
		},
	)
	if retry_err != nil {
		resp.Diagnostics.AddError(
			"Error deleting backend",
			"Could not delete backend, unexpected error: "+retry_err.Error(),
		)
		return
	}

}

func (r *backendResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	_, backendName, err := middleware.ResourceParseId(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Cannot parse import ID, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), backendName)...)
}
