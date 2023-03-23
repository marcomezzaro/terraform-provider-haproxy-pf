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
	_ resource.Resource                = &serverTemplateResource{}
	_ resource.ResourceWithConfigure   = &serverTemplateResource{}
	_ resource.ResourceWithImportState = &serverTemplateResource{}
)

// NewServerTemplateResource is a helper function to simplify the provider implementation.
func NewServerTemplateResource() resource.Resource {
	return &serverTemplateResource{}
}

// serverTemplateResource is the resource implementation.
type serverTemplateResource struct {
	client *middleware.Client
}

// serverTemplatesModel maps serverTemplates schema data.
type serverTemplateResourceModel struct {
	ID               types.String `tfsdk:"id"`
	Name             types.String `tfsdk:"name"`
	Fqdn             types.String `tfsdk:"fqdn"`
	Num_or_range     types.String `tfsdk:"num_or_range"`
	Port             types.Int64  `tfsdk:"port"`
	Prefix           types.String `tfsdk:"prefix"`
	Check            types.String `tfsdk:"check"`
	Resolvers        types.String `tfsdk:"resolvers"`
	ParentName       types.String `tfsdk:"parent_name"`

}

// Metadata returns the resource type name.
func (r *serverTemplateResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_server_template"
}

// Schema defines the schema for the resource.
func (r *serverTemplateResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Computed: true,
			},
			"fqdn": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"num_or_range": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"port": schema.Int64Attribute{
				Required: true,
				Optional: false,
			},
			"prefix": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"check": schema.StringAttribute{
				Required: true,
				Optional: false,
			},
			"resolvers": schema.StringAttribute{
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
func (r *serverTemplateResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*middleware.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *serverTemplateResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {

	// Retrieve values from plan
	var plan serverTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.ServerTemplate{
		Fqdn:             plan.Fqdn.ValueString(),
		Num_or_range:     plan.Num_or_range.ValueString(),
		Port:             plan.Port.ValueInt64(),
		Prefix:           plan.Prefix.ValueString(),
		Check: plan.Check.ValueString(),
		Resolvers: plan.Resolvers.ValueString(),
	}

	var response *models.ServerTemplate
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
			// Create new serverTemplate
			create_response, err := r.client.CreateServerTemplate(transaction.Id, payload, plan.ParentName.ValueString())
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
			"Error creating serverTemplate",
			"Could not create serverTemplate, unexpected error: "+retry_err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId(plan.ParentName.ValueString(), response.Prefix)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Prefix)
	plan.Fqdn = types.StringValue(response.Fqdn)
	plan.Num_or_range = types.StringValue(response.Num_or_range)
	plan.Port = types.Int64Value(response.Port)
	plan.Prefix = types.StringValue(response.Prefix)
	plan.Check = types.StringValue(response.Check)
	plan.Resolvers = types.StringValue(response.Resolvers)

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *serverTemplateResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	// Get current state
	var state serverTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, serverTemplateName, _ := middleware.ResourceParseId(ctx, state.ID.String())

	// Get refreshed serverTemplate
	response, err := r.client.GetServerTemplate(serverTemplateName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy ServerTemplate",
			"Could not read Haproxy ServerTemplate ID "+state.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Overwrite items with refreshed state
	id := middleware.CreateResourceId(parentName, response.Prefix)
	state.ID = types.StringValue(id)
	state.Name = types.StringValue(response.Prefix)
	state.Fqdn = types.StringValue(response.Fqdn)
	state.Num_or_range = types.StringValue(response.Num_or_range)
	state.Port = types.Int64Value(response.Port)
	state.Prefix = types.StringValue(response.Prefix)
	state.Check = types.StringValue(response.Check)
	state.Resolvers = types.StringValue(response.Resolvers)

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *serverTemplateResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	// retrieve values from state. extract ID
	var state serverTemplateResourceModel
	req.State.Get(ctx, &state)

	// Retrieve values from plan
	var plan serverTemplateResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// generate api request payload
	var payload = models.ServerTemplate{
		Fqdn:             plan.Fqdn.ValueString(),
		Num_or_range:     plan.Num_or_range.ValueString(),
		Port:             plan.Port.ValueInt64(),
		Prefix:           plan.Prefix.ValueString(),
		Check: plan.Check.ValueString(),
		Resolvers: plan.Resolvers.ValueString(),
	}
	parentName, serverTemplateName, err := middleware.ResourceParseId(ctx, state.ID.String())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting resource ID",
			"Could not update serverTemplate, unexpected error: "+err.Error(),
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
			// Update existing serverTemplate
			_, err = r.client.UpdateServerTemplate(transaction.Id, payload, parentName)
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
			"Error updating serverTemplate",
			"Could not update serverTemplate, unexpected error: "+retry_err.Error(),
		)
		return
	}

	response, err := r.client.GetServerTemplate(serverTemplateName, parentName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Haproxy ServerTemplate",
			"Could not read Haproxy ServerTemplate ID "+plan.ID.ValueString()+": "+err.Error(),
		)
		return
	}

	// Map response body to schema and populate Computed attribute values
	id := middleware.CreateResourceId(parentName, response.Prefix)
	plan.ID = types.StringValue(id)
	plan.Name = types.StringValue(response.Prefix)
	plan.Fqdn = types.StringValue(response.Fqdn)
	plan.Num_or_range = types.StringValue(response.Num_or_range)
	plan.Port = types.Int64Value(response.Port)
	plan.Prefix = types.StringValue(response.Prefix)
	plan.Check = types.StringValue(response.Check)
	plan.Resolvers = types.StringValue(response.Resolvers)


	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

// Delete deletes the resource and removes the Terraform state on success.
func (r *serverTemplateResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {

	// Retrieve values from state
	var state serverTemplateResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	parentName, serverTemplateName, _ := middleware.ResourceParseId(ctx, state.ID.String())

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
			// Delete existing serverTemplate
			err = r.client.DeleteServerTemplate(transaction.Id, serverTemplateName, parentName)
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
			"Error deleting serverTemplate",
			"Could not delete serverTemplate, unexpected error: "+retry_err.Error(),
		)
		return
	}
}

func (r *serverTemplateResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	parentName, serverTemplateName, err := middleware.ResourceParseId(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error importing resource",
			"Cannot parse import ID, unexpected error: "+err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("id"), req.ID)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("name"), serverTemplateName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("prefix"), serverTemplateName)...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("parent_name"), parentName)...)
}
