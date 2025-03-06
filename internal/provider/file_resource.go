package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/sashabaranov/go-openai"
)

var (
	_ resource.Resource                = &FileResource{}
	_ resource.ResourceWithImportState = &FileResource{}
)

func NewFileResource() resource.Resource {
	return &FileResource{}
}

type FileResource struct {
	client *openai.Client
}

type FileResourceModel struct {
	ID       types.String `tfsdk:"id"`
	FilePath types.String `tfsdk:"file_path"`
	FileName types.String `tfsdk:"filename"`
	Purpose  types.String `tfsdk:"purpose"`
	Bytes    types.Int64  `tfsdk:"bytes"`
}

func (r *FileResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_file"
}

func (r *FileResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: "Manages an OpenAI file.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description: "The ID of the file.",
				Computed:    true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"file_path": schema.StringAttribute{
				Description: "The path to the file to upload.",
				Required:    true,
			},
			"filename": schema.StringAttribute{
				Description: "The name of the file once uploaded to OpenAI.",
				Required:    true,
			},
			"purpose": schema.StringAttribute{
				Description: "The intended purpose of the file. Use 'fine-tune' for Fine-tuning, 'assistants' for Assistants and retrieval.",
				Required:    true,
			},
			"bytes": schema.Int64Attribute{
				Description: "The size of the file in bytes.",
				Computed:    true,
			},
		},
	}
}

func (r *FileResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*openai.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *openai.Client, got: %T. Please report this issue to the provider developers.", req.ProviderData),
		)
		return
	}

	r.client = client
}

func (r *FileResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var data FileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.CreateFile(ctx, openai.FileRequest{
		FilePath: data.FilePath.ValueString(),
		FileName: data.FileName.ValueString(),
		Purpose:  data.Purpose.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating File",
			fmt.Sprintf("Unable to create file, got error: %s", err),
		)
		return
	}

	data.ID = types.StringValue(file.ID)
	data.Bytes = types.Int64Value(int64(file.Bytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var data FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	file, err := r.client.GetFile(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading File",
			fmt.Sprintf("Unable to read file %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}

	data.ID = types.StringValue(file.ID)
	data.Bytes = types.Int64Value(int64(file.Bytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Files are immutable in OpenAI API, so we need to create a new one and delete the old one
	var data FileResourceModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create new file
	file, err := r.client.CreateFile(ctx, openai.FileRequest{
		FilePath: data.FilePath.ValueString(),
		FileName: data.FileName.ValueString(),
		Purpose:  data.Purpose.ValueString(),
	})
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating File",
			fmt.Sprintf("Unable to create file, got error: %s", err),
		)
		return
	}

	// Delete old file
	err = r.client.DeleteFile(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Error Deleting Old File",
			fmt.Sprintf("Unable to delete old file %s, got error: %s", state.ID.ValueString(), err),
		)
	}

	data.ID = types.StringValue(file.ID)
	data.Bytes = types.Int64Value(int64(file.Bytes))

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r *FileResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var data FileResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.client.DeleteFile(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting File",
			fmt.Sprintf("Unable to delete file %s, got error: %s", data.ID.ValueString(), err),
		)
		return
	}
}

// ImportState is called when importing a resource into Terraform.
// The import identifier is the OpenAI File ID.
func (r *FileResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve the file from OpenAI
	file, err := r.client.GetFile(ctx, req.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading File",
			fmt.Sprintf("Unable to read file %s, got error: %s", req.ID, err),
		)
		return
	}

	var data FileResourceModel
	data.ID = types.StringValue(file.ID)
	data.FileName = types.StringValue(file.FileName)
	data.Purpose = types.StringValue(file.Purpose)
	data.Bytes = types.Int64Value(int64(file.Bytes))
	// Since the actual file path is not available from the API,
	// we set it to a placeholder that the user must update
	data.FilePath = types.StringValue("REPLACE_WITH_LOCAL_FILE_PATH")

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
