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
var _ resource.Resource = &VectorStoreResource{}
var _ resource.ResourceWithImportState = &VectorStoreResource{}

func NewVectorStoreResource() resource.Resource {
	return &VectorStoreResource{}
}

// VectorStoreResource defines the resource implementation.
type VectorStoreResource struct {
	client *client.Client
}

// VectorStoreResourceModel describes the resource data model.
type VectorStoreResourceModel struct {
	ID           types.String `tfsdk:"id"`
	Name         types.String `tfsdk:"name"`
	CreatedAt    types.Int64  `tfsdk:"created_at"`
	Status       types.String `tfsdk:"status"`
	UsageBytes   types.Int64  `tfsdk:"usage_bytes"`
	ExpiresAt    types.Int64  `tfsdk:"expires_at"`
	ExpiresAfter *struct {
		Days   types.Int64  `tfsdk:"days"`
		Anchor types.String `tfsdk:"anchor"`
	} `tfsdk:"expires_after"`
	Metadata types.Map `tfsdk:"metadata"`
}

func (r *VectorStoreResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vector_store"
}

func (r *VectorStoreResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Manages an OpenAI vector store for efficient embedding storage and retrieval.",

		Blocks: map[string]schema.Block{
			"expires_after": schema.SingleNestedBlock{
				MarkdownDescription: "Configuration for vector store expiration.",
				Blocks:              map[string]schema.Block{},
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						MarkdownDescription: "Number of days after which the vector store expires.",
						Required:            true,
					},
					"anchor": schema.StringAttribute{
						MarkdownDescription: "Reference time for expiration calculation.",
						Optional:            true,
					},
				},
			},
		},

		Attributes: map[string]schema.Attribute{
			"metadata": schema.MapAttribute{
				ElementType:         types.StringType,
				Optional:            true,
				MarkdownDescription: "Metadata key-value pairs for the vector store.",
			},
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vector store.",
				Computed:            true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the vector store.",
				Required:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the vector store was created.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the vector store.",
				Computed:            true,
			},
			"usage_bytes": schema.Int64Attribute{
				MarkdownDescription: "The total size of the vector store in bytes.",
				Computed:            true,
			},
			"expires_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the vector store will expire.",
				Computed:            true,
			},
		},
	}
}

func (r *VectorStoreResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *VectorStoreResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan VectorStoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert metadata to map[string]interface{}
	metadata := make(map[string]interface{})
	if !plan.Metadata.IsNull() && !plan.Metadata.IsUnknown() {
		diags = plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the request
	createReq := openai.VectorStoreRequest{
		Name:     plan.Name.ValueString(),
		Metadata: metadata,
	}

	// Add expiration if configured
	if plan.ExpiresAfter != nil {
		createReq.ExpiresAfter = &openai.VectorStoreExpires{
			Days: int(plan.ExpiresAfter.Days.ValueInt64()),
		}
		if !plan.ExpiresAfter.Anchor.IsNull() {
			createReq.ExpiresAfter.Anchor = plan.ExpiresAfter.Anchor.ValueString()
		}
	}

	tflog.Debug(ctx, "Creating vector store", map[string]interface{}{
		"name": createReq.Name,
	})

	result, err := r.client.OpenAI.CreateVectorStore(ctx, createReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Vector Store",
			fmt.Sprintf("Unable to create vector store: %s", r.client.HandleError(err)),
		)
		return
	}

	// Map response to model
	plan.ID = types.StringValue(result.ID)
	plan.CreatedAt = types.Int64Value(result.CreatedAt)
	plan.Status = types.StringValue(result.Status)
	plan.UsageBytes = types.Int64Value(int64(result.UsageBytes))
	if result.ExpiresAt != nil {
		plan.ExpiresAt = types.Int64Value(int64(*result.ExpiresAt))
	}

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state VectorStoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := r.client.OpenAI.RetrieveVectorStore(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vector Store",
			fmt.Sprintf("Unable to read vector store: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state
	state.Name = types.StringValue(result.Name)
	state.CreatedAt = types.Int64Value(result.CreatedAt)
	state.Status = types.StringValue(result.Status)
	state.UsageBytes = types.Int64Value(int64(result.UsageBytes))
	if result.ExpiresAt != nil {
		state.ExpiresAt = types.Int64Value(int64(*result.ExpiresAt))
	}

	// Save updated state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan VectorStoreResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Convert metadata to map[string]interface{}
	metadata := make(map[string]interface{})
	if !plan.Metadata.IsNull() && !plan.Metadata.IsUnknown() {
		diags = plan.Metadata.ElementsAs(ctx, &metadata, false)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

	// Build the request
	updateReq := openai.VectorStoreRequest{
		Name:     plan.Name.ValueString(),
		Metadata: metadata,
	}

	// Add expiration if configured
	if plan.ExpiresAfter != nil {
		updateReq.ExpiresAfter = &openai.VectorStoreExpires{
			Days: int(plan.ExpiresAfter.Days.ValueInt64()),
		}
		if !plan.ExpiresAfter.Anchor.IsNull() {
			updateReq.ExpiresAfter.Anchor = plan.ExpiresAfter.Anchor.ValueString()
		}
	}

	tflog.Debug(ctx, "Updating vector store", map[string]interface{}{
		"id": plan.ID.ValueString(),
	})

	result, err := r.client.OpenAI.ModifyVectorStore(ctx, plan.ID.ValueString(), updateReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Vector Store",
			fmt.Sprintf("Unable to update vector store: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update state with response
	plan.CreatedAt = types.Int64Value(result.CreatedAt)
	plan.Status = types.StringValue(result.Status)
	plan.UsageBytes = types.Int64Value(int64(result.UsageBytes))
	if result.ExpiresAt != nil {
		plan.ExpiresAt = types.Int64Value(int64(*result.ExpiresAt))
	}

	// Save updated state
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *VectorStoreResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state VectorStoreResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	_, err := r.client.OpenAI.DeleteVectorStore(ctx, state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Vector Store",
			fmt.Sprintf("Unable to delete vector store: %s", r.client.HandleError(err)),
		)
		return
	}
}

func (r *VectorStoreResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
