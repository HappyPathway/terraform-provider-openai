package resources

import (
	"context"
	"fmt"
	"strings"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
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
var _ resource.Resource = &EmbeddingResource{}
var _ resource.ResourceWithImportState = &EmbeddingResource{}

func NewEmbeddingResource() resource.Resource {
	return &EmbeddingResource{}
}

// EmbeddingResource defines the resource implementation.
type EmbeddingResource struct {
	client *client.Client
}

// EmbeddingResourceModel describes the resource data model.
type EmbeddingResourceModel struct {
	ID        types.String `tfsdk:"id"`
	Model     types.String `tfsdk:"model"`
	Input     types.String `tfsdk:"input"`
	Embedding types.List   `tfsdk:"embedding"`
}

func (r *EmbeddingResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_embedding"
}

func (r *EmbeddingResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Generate vector embeddings from text using OpenAI's embedding models.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "ID of the model to use for generating embeddings (e.g., 'text-embedding-ada-002').",
				Required:            true,
			},
			"input": schema.StringAttribute{
				MarkdownDescription: "The text to generate embeddings for.",
				Required:            true,
			},
			"embedding": schema.ListAttribute{
				ElementType:         types.Float64Type,
				MarkdownDescription: "The vector embedding representing the input text.",
				Computed:            true,
			},
		},
	}
}

func (r *EmbeddingResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *EmbeddingResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan EmbeddingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the embedding request
	embeddingReq := openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(plan.Model.ValueString()),
		Input: plan.Input.ValueString(),
	}

	tflog.Debug(ctx, "Creating embedding", map[string]interface{}{
		"model": embeddingReq.Model,
	})

	// Call OpenAI API
	result, err := r.client.OpenAI.CreateEmbeddings(ctx, embeddingReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Embedding",
			fmt.Sprintf("Unable to create embedding: %s", r.client.HandleError(err)),
		)
		return
	}

	if len(result.Data) == 0 {
		resp.Diagnostics.AddError(
			"Error Creating Embedding",
			"OpenAI API returned an empty embedding result",
		)
		return
	}

	// Convert the embedding vector to a list of values
	embedding := make([]attr.Value, len(result.Data[0].Embedding))
	for i, value := range result.Data[0].Embedding {
		embedding[i] = types.Float64Value(float64(value))
	}

	embeddingList, diags := types.ListValue(types.Float64Type, embedding)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Generate a deterministic ID based on model and input
	inputHash := createHashFromInput(plan.Input.ValueString())
	plan.ID = types.StringValue(fmt.Sprintf("embed-%s-%s", plan.Model.ValueString(), inputHash))
	plan.Embedding = embeddingList

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EmbeddingResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state EmbeddingResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Embedding is stateless - we just keep whatever is in the state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *EmbeddingResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan EmbeddingResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the embedding request
	embeddingReq := openai.EmbeddingRequest{
		Model: openai.EmbeddingModel(plan.Model.ValueString()),
		Input: plan.Input.ValueString(),
	}

	tflog.Debug(ctx, "Updating embedding", map[string]interface{}{
		"model": embeddingReq.Model,
	})

	// Call OpenAI API
	result, err := r.client.OpenAI.CreateEmbeddings(ctx, embeddingReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Embedding",
			fmt.Sprintf("Unable to create embedding: %s", r.client.HandleError(err)),
		)
		return
	}

	if len(result.Data) == 0 {
		resp.Diagnostics.AddError(
			"Error Creating Embedding",
			"OpenAI API returned an empty embedding result",
		)
		return
	}

	// Convert the embedding vector to a list of values
	embedding := make([]attr.Value, len(result.Data[0].Embedding))
	for i, value := range result.Data[0].Embedding {
		embedding[i] = types.Float64Value(float64(value))
	}

	embeddingList, diags := types.ListValue(types.Float64Type, embedding)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Update the ID if the input or model has changed
	inputHash := createHashFromInput(plan.Input.ValueString())
	plan.ID = types.StringValue(fmt.Sprintf("embed-%s-%s", plan.Model.ValueString(), inputHash))
	plan.Embedding = embeddingList

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *EmbeddingResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// No API call needed for deletion since embeddings are stateless
	tflog.Info(ctx, "Deleted embedding resource", map[string]interface{}{
		"id": req.State.GetAttribute(ctx, path.Root("id"), nil),
	})
}

func (r *EmbeddingResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Helper function to create a simple hash from the input text
func createHashFromInput(input string) string {
	// For simplicity, just take first 8 chars and sanitize
	if len(input) > 8 {
		input = input[:8]
	}
	input = strings.ReplaceAll(input, " ", "")
	input = strings.ReplaceAll(input, "\n", "")
	input = strings.ReplaceAll(input, "\t", "")
	return input
}
