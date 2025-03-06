package resources

import (
	"context"
	"fmt"

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
var _ resource.Resource = &FileResource{}
var _ resource.ResourceWithImportState = &FileResource{}

func NewFileResource() resource.Resource {
	return &FileResource{}
}

// FileResource defines the resource implementation.
type FileResource struct {
	client *client.Client
}

// FileResourceModel describes the resource data model.
type FileResourceModel struct {
	ID            types.String `tfsdk:"id"`
	Filename      types.String `tfsdk:"filename"`
	FilePath      types.String `tfsdk:"file_path"`
	Purpose       types.String `tfsdk:"purpose"`
	ObjectID      types.String `tfsdk:"object_id"`
	Bytes         types.Int64  `tfsdk:"bytes"`
	CreatedAt     types.Int64  `tfsdk:"created_at"`
	Status        types.String `tfsdk:"status"`
	StatusDetails types.String `tfsdk:"status_details"`
}

func (r *FileResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *FileResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Upload and manage files on the OpenAI platform for fine-tuning or use with assistants.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"filename": schema.StringAttribute{
				MarkdownDescription: "The name of the file being uploaded.",
				Required:            true,
			},
			"file_path": schema.StringAttribute{
				MarkdownDescription: "The local path to the file to upload.",
				Required:            true,
			},
			"purpose": schema.StringAttribute{
				MarkdownDescription: "The purpose of the file. Allowed values are 'fine-tune', 'fine-tune-results', or 'assistants'.",
				Required:            true,
			},
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The OpenAI ID of the uploaded file.",
				Computed:            true,
			},
			"bytes": schema.Int64Attribute{
				MarkdownDescription: "The size of the file in bytes.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the file was created.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the file. Can be 'uploaded', 'processed', or 'error'.",
				Computed:            true,
			},
			"status_details": schema.StringAttribute{
				MarkdownDescription: "Additional details about the file's status, particularly useful for errors.",
				Computed:            true,
			},
		},
	}
}

func (r *FileResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FileResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create file request
	fileReq := openai.FileRequest{
		FilePath: plan.FilePath.ValueString(),
		Purpose:  plan.Purpose.ValueString(),
	}

	tflog.Debug(ctx, "Creating file", map[string]interface{}{
		"filename": plan.Filename.ValueString(),
		"purpose":  plan.Purpose.ValueString(),
	})

	// Upload the file - no need to set headers manually as they are handled by the transport
	file, err := r.client.OpenAI.CreateFile(ctx, fileReq)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating File",
			fmt.Sprintf("Unable to create file: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	plan.ID = types.StringValue(file.ID)
	plan.ObjectID = types.StringValue(file.ID)
	plan.Filename = types.StringValue(file.FileName)
	plan.Bytes = types.Int64Value(int64(file.Bytes))
	plan.CreatedAt = types.Int64Value(int64(file.CreatedAt))
	plan.Status = types.StringValue(file.Status)
	plan.StatusDetails = types.StringValue(file.StatusDetails)

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileID := state.ObjectID.ValueString()
	if fileID == "" {
		resp.Diagnostics.AddError(
			"Error Reading File",
			"File ID is empty. Cannot retrieve file details.",
		)
		return
	}

	tflog.Debug(ctx, "Reading file", map[string]interface{}{
		"file_id": fileID,
	})

	// Retrieve file information
	file, err := r.client.OpenAI.GetFile(ctx, fileID)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			// File doesn't exist anymore, remove from state
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading File",
			fmt.Sprintf("Unable to read file details: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the computed attributes in state
	state.ObjectID = types.StringValue(file.ID)
	state.Bytes = types.Int64Value(int64(file.Bytes))
	state.CreatedAt = types.Int64Value(int64(file.CreatedAt))
	state.Status = types.StringValue(file.Status)
	state.StatusDetails = types.StringValue(file.StatusDetails)

	// Save into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan FileResourceModel
	var state FileResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete the existing file
	fileID := state.ObjectID.ValueString()
	if fileID != "" {
		err := r.client.OpenAI.DeleteFile(ctx, fileID)
		if err != nil {
			// If file doesn't exist, continue with creation
			if apiErr, ok := err.(*openai.APIError); !ok || apiErr.HTTPStatusCode != 404 {
				resp.Diagnostics.AddError(
					"Error Deleting File",
					fmt.Sprintf("Unable to delete file before recreation: %s", r.client.HandleError(err)),
				)
				return
			}
		}
	}

	// Create file request - no need to set headers manually
	file, err := r.client.OpenAI.CreateFile(ctx, openai.FileRequest{
		FilePath: plan.FilePath.ValueString(),
		Purpose:  plan.Purpose.ValueString(),
	})

	tflog.Debug(ctx, "Updating file (recreate)", map[string]interface{}{
		"filename": plan.Filename.ValueString(),
		"purpose":  plan.Purpose.ValueString(),
	})

	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating File",
			fmt.Sprintf("Unable to create file: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	plan.ID = types.StringValue(file.ID)
	plan.ObjectID = types.StringValue(file.ID)
	plan.Bytes = types.Int64Value(int64(file.Bytes))
	plan.CreatedAt = types.Int64Value(int64(file.CreatedAt))
	plan.Status = types.StringValue(file.Status)
	plan.StatusDetails = types.StringValue(file.StatusDetails)

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fileID := state.ObjectID.ValueString()
	if fileID == "" {
		// Nothing to delete
		return
	}

	tflog.Debug(ctx, "Deleting file", map[string]interface{}{
		"file_id": fileID,
	})

	err := r.client.OpenAI.DeleteFile(ctx, fileID)
	if err != nil {
		// If file doesn't exist, don't return an error
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(
			"Error Deleting File",
			fmt.Sprintf("Unable to delete file: %s", r.client.HandleError(err)),
		)
		return
	}
}

func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("object_id"), req, resp)
}
