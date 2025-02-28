package openai

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
	"github.com/openai/openai-go/shared"
)

func resourceOpenAIFineTuningJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIFineTuningJobCreate,
		ReadContext:   resourceOpenAIFineTuningJobRead,
		DeleteContext: resourceOpenAIFineTuningJobDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},

		Schema: map[string]*schema.Schema{
			"model": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the base model to fine-tune",
				DiffSuppressFunc: func(k, old, new string, d *schema.ResourceData) bool {
					// Suppress differences between base model and specific versions
					// e.g., "gpt-3.5-turbo" vs "gpt-3.5-turbo-0125"
					return strings.HasPrefix(old, new) || strings.HasPrefix(new, old)
				},
			},
			"training_file": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"hyperparameters": {
				Type:     schema.TypeList,
				Optional: true,
				ForceNew: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"batch_size": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
						"learning_rate_multiplier": {
							Type:     schema.TypeFloat,
							Optional: true,
							ForceNew: true,
						},
						"n_epochs": {
							Type:     schema.TypeInt,
							Optional: true,
							ForceNew: true,
						},
					},
				},
			},
			"validation_file": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"suffix": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"fine_tuned_model": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"trained_tokens": {
				Type:     schema.TypeInt,
				Computed: true,
			},
			"result_files": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp for when the fine-tuning job was created",
			},
			"finished_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp for when the fine-tuning job was finished",
			},
			"error": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"code": {
							Type:     schema.TypeString,
							Computed: true,
						},
						"message": {
							Type:     schema.TypeString,
							Computed: true,
						},
					},
				},
				Description: "Error details if the fine-tuning job failed",
			},
		},
	}
}

func resourceOpenAIFineTuningJobCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	params := openaiapi.FineTuningJobNewParams{
		Model:        openaiapi.F(openaiapi.FineTuningJobNewParamsModel(d.Get("model").(string))),
		TrainingFile: openaiapi.F(d.Get("training_file").(string)),
	}

	if v, ok := d.GetOk("validation_file"); ok {
		params.ValidationFile = openaiapi.F(v.(string))
	}

	if v, ok := d.GetOk("suffix"); ok {
		params.Suffix = openaiapi.F(v.(string))
	}

	if v, ok := d.GetOk("hyperparameters"); ok {
		hyperList := v.([]interface{})
		if len(hyperList) > 0 {
			hyperMap := hyperList[0].(map[string]interface{})
			hyperParams := openaiapi.FineTuningJobNewParamsHyperparameters{}

			if batchSize, ok := hyperMap["batch_size"]; ok {
				hyperParams.BatchSize = openaiapi.F[openaiapi.FineTuningJobNewParamsHyperparametersBatchSizeUnion](shared.UnionInt(batchSize.(int)))
			}

			if lr, ok := hyperMap["learning_rate_multiplier"]; ok {
				hyperParams.LearningRateMultiplier = openaiapi.F[openaiapi.FineTuningJobNewParamsHyperparametersLearningRateMultiplierUnion](shared.UnionFloat(lr.(float64)))
			}

			if epochs, ok := hyperMap["n_epochs"]; ok {
				hyperParams.NEpochs = openaiapi.F[openaiapi.FineTuningJobNewParamsHyperparametersNEpochsUnion](shared.UnionInt(epochs.(int)))
			}

			params.Hyperparameters = openaiapi.F(hyperParams)
		}
	}

	fineTuningJob, err := client.FineTuning.Jobs.New(ctx, params)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fineTuningJob.ID)
	return resourceOpenAIFineTuningJobRead(ctx, d, m)
}

func resourceOpenAIFineTuningJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	fineTuningJob, err := client.FineTuning.Jobs.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("model", fineTuningJob.Model)
	d.Set("training_file", fineTuningJob.TrainingFile)
	d.Set("validation_file", fineTuningJob.ValidationFile)
	d.Set("status", fineTuningJob.Status)
	d.Set("fine_tuned_model", fineTuningJob.FineTunedModel)
	d.Set("trained_tokens", fineTuningJob.TrainedTokens)
	d.Set("result_files", fineTuningJob.ResultFiles)
	d.Set("created_at", fineTuningJob.CreatedAt)
	d.Set("finished_at", fineTuningJob.FinishedAt)

	// The Error field is a required field that can be null in the API response
	// We need to check if it has actual error data before setting it
	if fineTuningJob.Error.Code != "" || fineTuningJob.Error.Message != "" {
		error := []interface{}{
			map[string]interface{}{
				"code":    fineTuningJob.Error.Code,
				"message": fineTuningJob.Error.Message,
			},
		}
		d.Set("error", error)
	}

	return nil
}

func resourceOpenAIFineTuningJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	// Since we cannot delete a fine-tuning job, we'll try to cancel it if it's still running
	fineTuningJob, err := client.FineTuning.Jobs.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	if fineTuningJob.Status == openaiapi.FineTuningJobStatusQueued ||
		fineTuningJob.Status == openaiapi.FineTuningJobStatusRunning {
		_, err = client.FineTuning.Jobs.Cancel(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// Remove from state
	d.SetId("")
	return nil
}
