package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/sashabaranov/go-openai"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ resource.Resource = &VectorStoreFileResource{}
var _ resource.ResourceWithImportState = &VectorStoreFileResource{}

func NewVectorStoreFileResource() resource.Resource {
	return &VectorStoreFileResource{}
}

// VectorStoreFileResource defines the resource implementation.
type VectorStoreFileResource struct {
	client *client.Client
}

// VectorStoreFileResourceModel describes the resource data model.
type VectorStoreFileResourceModel struct {
	ID            types.String `tfsdk:"id"`
	VectorStoreID types.String `tfsdk:"vector_store_id"`
	FileID        types.String `tfsdk:"file_id"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	UsageBytes    types.Int64  `tfsdk:"usage_bytes"`
	Status        types.String `tfsdk:"status"`
}

func (r *VectorStoreFileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vector_store_file"
}

func (r *VectorStoreFileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages a file within an OpenAI vector store.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The unique identifier for this vector store file.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"vector_store_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vector store to add the file to.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"file_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the file to add to the vector store.",
				Required:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the file was added to the vector store.",
				Computed:            true,
			},
			"usage_bytes": schema.Int64Attribute{
				MarkdownDescription: "The size of the file in bytes.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the file in the vector store.",
				Computed:            true,
			},
		},
	}
}

func (r *VectorStoreFileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *VectorStoreFileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VectorStoreFileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Debug(ctx, "Creating vector store file", map[string]interface{}{
		"vector_store_id": plan.VectorStoreID.ValueString(),
		"file_id":         plan.FileID.ValueString(),
	})

	result, err := r.client.OpenAI.CreateVectorStoreFile(ctx, plan.VectorStoreID.ValueString(), openai.VectorStoreFileRequest{
		FileID: plan.FileID.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Vector Store File",
			fmt.Sprintf("Unable to add file to vector store: %s", r.client.HandleError(err)),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(result.ID)
	plan.CreatedAt = types.Int64Value(result.CreatedAt)
	plan.Status = types.StringValue(result.Status)
	plan.UsageBytes = types.Int64Value(int64(result.UsageBytes))

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreFileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VectorStoreFileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.OpenAI.RetrieveVectorStoreFile(ctx, state.VectorStoreID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vector Store File",
			fmt.Sprintf("Unable to read vector store file: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state
	state.CreatedAt = types.Int64Value(result.CreatedAt)
	state.Status = types.StringValue(result.Status)
	state.UsageBytes = types.Int64Value(int64(result.UsageBytes))

	// Save updated state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreFileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Files in vector stores can't be updated, they need to be replaced
	var plan VectorStoreFileResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VectorStoreFileResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Since the file can't be updated, we just read the current state
	result, err := r.client.OpenAI.RetrieveVectorStoreFile(ctx, state.VectorStoreID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vector Store File",
			fmt.Sprintf("Unable to read vector store file: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state
	plan.CreatedAt = types.Int64Value(result.CreatedAt)
	plan.Status = types.StringValue(result.Status)
	plan.UsageBytes = types.Int64Value(int64(result.UsageBytes))

	// Save updated state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreFileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VectorStoreFileResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.OpenAI.DeleteVectorStoreFile(ctx, state.VectorStoreID.ValueString(), state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Vector Store File",
			fmt.Sprintf("Unable to delete vector store file: %s", r.client.HandleError(err)),
		)
		return
	}
}

func (r *VectorStoreFileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Expected import ID format: vector_store_id:file_id
	idParts := strings.Split(req.ID, ":")
	if len(idParts) != 2 || idParts[0] == "" || idParts[1] == "" {
		resp.Diagnostics.AddError(
			"Invalid Import ID",
			`The import ID must be in the format "vector_store_id:file_id"`,
		)
		return
	}

	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("vector_store_id"), idParts[0])...)
	resp.Diagnostics.Append(resp.State.SetAttribute(ctx, path.Root("file_id"), idParts[1])...)
}
