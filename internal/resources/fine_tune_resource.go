package resources

import (
	"context"
	"fmt"
	"time"

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
var _ resource.Resource = &FineTuneResource{}
var _ resource.ResourceWithImportState = &FineTuneResource{}

func NewFineTuneResource() resource.Resource {
	return &FineTuneResource{}
}

// FineTuneResource defines the resource implementation.
type FineTuneResource struct {
	client *client.Client
}

// FineTuneResourceModel describes the resource data model.
type FineTuneResourceModel struct {
	ID                           types.String  `tfsdk:"id"`
	TrainingFileID               types.String  `tfsdk:"training_file_id"`
	ValidationFileID             types.String  `tfsdk:"validation_file_id"`
	Model                        types.String  `tfsdk:"model"`
	Epochs                       types.Int64   `tfsdk:"epochs"`
	BatchSize                    types.Int64   `tfsdk:"batch_size"`
	LearningRateMultiplier       types.Float64 `tfsdk:"learning_rate_multiplier"`
	PromptLossWeight             types.Float64 `tfsdk:"prompt_loss_weight"`
	ComputeClassificationMetrics types.Bool    `tfsdk:"compute_classification_metrics"`
	ClassificationNClasses       types.Int64   `tfsdk:"classification_n_classes"`
	ClassificationPositiveClass  types.String  `tfsdk:"classification_positive_class"`
	Suffix                       types.String  `tfsdk:"suffix"`

	// Computed outputs
	ObjectID       types.String `tfsdk:"object_id"`
	Status         types.String `tfsdk:"status"`
	FineTunedModel types.String `tfsdk:"fine_tuned_model"`
	OrganizationID types.String `tfsdk:"organization_id"`
	ResultFiles    types.List   `tfsdk:"result_files"`
	CreatedAt      types.Int64  `tfsdk:"created_at"`
}

func (r *FineTuneResource) Metadata(ctx context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_fine_tune"
}

func (r *FineTuneResource) Schema(ctx context.Context, req resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Create and manage fine-tuning jobs to customize OpenAI models for specific use cases.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				MarkdownDescription: "Unique identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"training_file_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the file to use for training the model. You can upload files using the `openai_file` resource.",
				Required:            true,
			},
			"validation_file_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the file to use for validation of the fine-tuned model.",
				Optional:            true,
			},
			"model": schema.StringAttribute{
				MarkdownDescription: "The base model to fine-tune (e.g., 'gpt-3.5-turbo', 'davinci', etc).",
				Required:            true,
			},
			"epochs": schema.Int64Attribute{
				MarkdownDescription: "The number of epochs to train the model for.",
				Optional:            true,
			},
			"batch_size": schema.Int64Attribute{
				MarkdownDescription: "The batch size to use for training.",
				Optional:            true,
			},
			"learning_rate_multiplier": schema.Float64Attribute{
				MarkdownDescription: "The learning rate multiplier to use for training.",
				Optional:            true,
			},
			"prompt_loss_weight": schema.Float64Attribute{
				MarkdownDescription: "The weight to use for prompt tokens loss.",
				Optional:            true,
			},
			"compute_classification_metrics": schema.BoolAttribute{
				MarkdownDescription: "If true, compute classification metrics like accuracy, F1 score, etc.",
				Optional:            true,
			},
			"classification_n_classes": schema.Int64Attribute{
				MarkdownDescription: "The number of classes for multiclass classification.",
				Optional:            true,
			},
			"classification_positive_class": schema.StringAttribute{
				MarkdownDescription: "The positive class for binary classification.",
				Optional:            true,
			},
			"suffix": schema.StringAttribute{
				MarkdownDescription: "A suffix to append to the fine-tuned model name.",
				Optional:            true,
			},

			// Computed outputs
			"object_id": schema.StringAttribute{
				MarkdownDescription: "The ID of the fine-tuning job.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The status of the fine-tuning job (e.g., 'pending', 'running', 'succeeded', 'failed').",
				Computed:            true,
			},
			"fine_tuned_model": schema.StringAttribute{
				MarkdownDescription: "The name of the fine-tuned model once training is complete.",
				Computed:            true,
			},
			"organization_id": schema.StringAttribute{
				MarkdownDescription: "The organization that owns the fine-tuned model.",
				Computed:            true,
			},
			"result_files": schema.ListAttribute{
				ElementType:         types.StringType,
				MarkdownDescription: "File IDs of any result files (e.g., validation results).",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The Unix timestamp (in seconds) for when the fine-tuning job was created.",
				Computed:            true,
			},
		},
	}
}

func (r *FineTuneResource) Configure(ctx context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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

func (r *FineTuneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan FineTuneResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create the fine-tuning request using the latest API
	fineTuneReq := openai.FineTuningJobRequest{
		TrainingFile: plan.TrainingFileID.ValueString(),
		Model:        plan.Model.ValueString(),
	}

	// Set optional parameters
	if !plan.ValidationFileID.IsNull() {
		fineTuneReq.ValidationFile = plan.ValidationFileID.ValueString()
	}

	if !plan.Epochs.IsNull() {
		fineTuneReq.Hyperparameters = &openai.Hyperparameters{
			NEpochs: int(plan.Epochs.ValueInt64()),
		}
	}

	if !plan.BatchSize.IsNull() {
		// Note: BatchSize is no longer directly supported in the new API
		// We'll log a warning but not fail
		tflog.Warn(ctx, "batch_size parameter is no longer directly supported in the OpenAI API and will be ignored", map[string]interface{}{
			"batch_size": plan.BatchSize.ValueInt64(),
		})
	}

	if !plan.LearningRateMultiplier.IsNull() {
		// Note: LearningRateMultiplier is no longer directly supported in the new API
		// We'll log a warning but not fail
		tflog.Warn(ctx, "learning_rate_multiplier parameter is no longer directly supported in the OpenAI API and will be ignored", map[string]interface{}{
			"learning_rate_multiplier": plan.LearningRateMultiplier.ValueFloat64(),
		})
	}

	if !plan.Suffix.IsNull() {
		fineTuneReq.Suffix = plan.Suffix.ValueString()
	}

	tflog.Debug(ctx, "Creating fine-tune job", map[string]interface{}{
		"model":         fineTuneReq.Model,
		"training_file": fineTuneReq.TrainingFile,
	})

	// Submit the fine-tuning job
	fineTune, err := r.client.OpenAI.CreateFineTuningJob(ctx, fineTuneReq)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Creating Fine-Tune Job",
			fmt.Sprintf("Unable to create fine-tuning job: %s", r.client.HandleError(err)),
		)
		return
	}

	// Wait for job to be created
	time.Sleep(2 * time.Second)

	// Poll for job status
	fineTune, err = r.client.OpenAI.RetrieveFineTuningJob(ctx, fineTune.ID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Retrieving Fine-Tune Job",
			fmt.Sprintf("Unable to retrieve fine-tuning job status: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state with the returned data
	plan.ID = types.StringValue(fineTune.ID)
	plan.ObjectID = types.StringValue(fineTune.ID)
	plan.Status = types.StringValue(fineTune.Status)
	plan.FineTunedModel = types.StringValue(fineTune.FineTunedModel)
	plan.OrganizationID = types.StringValue(fineTune.OrganizationID)
	plan.CreatedAt = types.Int64Value(int64(fineTune.CreatedAt))

	// Convert result files to list if present
	resultFiles := []string{}
	for _, file := range fineTune.ResultFiles {
		resultFiles = append(resultFiles, file)
	}

	resultFilesList, diags := types.ListValueFrom(ctx, types.StringType, resultFiles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	plan.ResultFiles = resultFilesList

	// Save into state
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *FineTuneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state FineTuneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fineTuneID := state.ObjectID.ValueString()
	if fineTuneID == "" {
		resp.Diagnostics.AddError(
			"Error Reading Fine-Tune Job",
			"Fine-tune ID is empty. Cannot retrieve fine-tune details.",
		)
		return
	}

	tflog.Debug(ctx, "Reading fine-tune job", map[string]interface{}{
		"fine_tune_id": fineTuneID,
	})

	// Retrieve fine-tune information
	fineTune, err := r.client.OpenAI.RetrieveFineTuningJob(ctx, fineTuneID)
	if err != nil {
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			// Fine-tune doesn't exist anymore, remove from state
			resp.State.RemoveResource(ctx)
			return
		}

		resp.Diagnostics.AddError(
			"Error Reading Fine-Tune Job",
			fmt.Sprintf("Unable to read fine-tune details: %s", r.client.HandleError(err)),
		)
		return
	}

	// Update the state
	state.Status = types.StringValue(fineTune.Status)
	state.FineTunedModel = types.StringValue(fineTune.FineTunedModel)
	state.OrganizationID = types.StringValue(fineTune.OrganizationID)

	// Convert result files to list if present
	resultFiles := []string{}
	for _, file := range fineTune.ResultFiles {
		resultFiles = append(resultFiles, file)
	}

	resultFilesList, diags := types.ListValueFrom(ctx, types.StringType, resultFiles)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	state.ResultFiles = resultFilesList

	// Save into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &state)...)
}

func (r *FineTuneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Fine-tunes cannot be updated after creation
	resp.Diagnostics.AddError(
		"Error Updating Fine-Tune",
		"OpenAI fine-tune jobs cannot be updated after creation. Delete and recreate the resource to change parameters.",
	)
}

func (r *FineTuneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state FineTuneResourceModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	fineTuneID := state.ObjectID.ValueString()
	if fineTuneID == "" {
		// Nothing to delete
		return
	}

	tflog.Debug(ctx, "Cancelling fine-tune job", map[string]interface{}{
		"fine_tune_id": fineTuneID,
	})

	// Cancel the fine-tuning job if it's still in progress
	fineTune, err := r.client.OpenAI.RetrieveFineTuningJob(ctx, fineTuneID)
	if err != nil {
		// If fine-tune doesn't exist, don't return an error
		if apiErr, ok := err.(*openai.APIError); ok && apiErr.HTTPStatusCode == 404 {
			return
		}

		resp.Diagnostics.AddError(
			"Error Retrieving Fine-Tune Job",
			fmt.Sprintf("Unable to retrieve fine-tuning job: %s", r.client.HandleError(err)),
		)
		return
	}

	if fineTune.Status == "pending" || fineTune.Status == "running" {
		_, err = r.client.OpenAI.CancelFineTuningJob(ctx, fineTuneID)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Cancelling Fine-Tune Job",
				fmt.Sprintf("Unable to cancel fine-tuning job: %s", r.client.HandleError(err)),
			)
			return
		}
	}

	// Note: We can't delete fine-tuned models via the API, just cancel the job
	tflog.Info(ctx, "Fine-tune job cancelled, but the fine-tuned model (if any) still exists in OpenAI", map[string]interface{}{
		"fine_tune_id": fineTuneID,
		"model":        state.FineTunedModel.ValueString(),
	})
}

func (r *FineTuneResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("object_id"), req, resp)
}
