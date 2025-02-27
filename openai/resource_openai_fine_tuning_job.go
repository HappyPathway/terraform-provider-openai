package openai

import (
	"context"
	"strings"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenAIFineTuningJob() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIFineTuningJobCreate,
		ReadContext:   resourceOpenAIFineTuningJobRead,
		DeleteContext: resourceOpenAIFineTuningJobDelete,

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
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The ID of the training file",
			},
			"validation_file": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "The ID of the validation file",
			},
			"n_epochs": {
				Type:        schema.TypeInt,
				Optional:    true,
				ForceNew:    true,
				Description: "The number of epochs to train the model for",
			},
			"suffix": {
				Type:        schema.TypeString,
				Optional:    true,
				ForceNew:    true,
				Description: "A string of up to 40 characters that will be added to your fine-tuned model name",
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
			"status": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The status of the fine-tuning job",
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
	client := m.(*Client)

	req := &CreateFineTuningJobRequest{
		Model:          d.Get("model").(string),
		TrainingFile:   d.Get("training_file").(string),
		ValidationFile: d.Get("validation_file").(string),
		Suffix:         d.Get("suffix").(string),
	}

	if v, ok := d.GetOk("n_epochs"); ok {
		req.Hyperparameters.NEpochs = v.(int)
	}

	job, err := client.CreateFineTuningJob(ctx, req)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(job.ID)

	return resourceOpenAIFineTuningJobRead(ctx, d, m)
}

func resourceOpenAIFineTuningJobRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	job, err := client.GetFineTuningJob(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("model", job.Model)
	d.Set("training_file", job.TrainingFile)
	d.Set("validation_file", job.ValidationFile)
	d.Set("status", job.Status)
	d.Set("created_at", job.CreatedAt)
	d.Set("finished_at", job.FinishedAt)

	if job.Error != nil {
		error := []interface{}{
			map[string]interface{}{
				"code":    job.Error.Code,
				"message": job.Error.Message,
			},
		}
		d.Set("error", error)
	}

	return nil
}

func resourceOpenAIFineTuningJobDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	// Try to cancel the job if it's still running
	if d.Get("status").(string) == "running" {
		_, err := client.CancelFineTuningJob(ctx, d.Id())
		if err != nil {
			return diag.FromErr(err)
		}
	}

	// The job will remain in the API but we remove it from state
	d.SetId("")
	return nil
}
