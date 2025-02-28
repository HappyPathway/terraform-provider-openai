package openai

import (
	"context"
	"io"
	"os"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	openaiapi "github.com/openai/openai-go"
)

func resourceOpenAIFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIFileCreate,
		ReadContext:   resourceOpenAIFileRead,
		DeleteContext: resourceOpenAIFileDelete,
		Schema: map[string]*schema.Schema{
			"file": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "Path to the file to upload",
			},
			"purpose": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The intended purpose of the file. Possible values are: 'fine-tune', 'fine-tune-results', 'assistants', or 'assistants_output'",
			},
			"bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the file in bytes",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Creation timestamp",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the file",
			},
		},
	}
}

func resourceOpenAIFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	filePath := d.Get("file").(string)
	file, err := os.Open(filePath)
	if err != nil {
		return diag.FromErr(err)
	}
	defer file.Close()

	fileObj, err := client.Files.New(ctx, openaiapi.FileNewParams{
		File:    openaiapi.F[io.Reader](file),
		Purpose: openaiapi.F(openaiapi.FilePurpose(d.Get("purpose").(string))),
	})
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(fileObj.ID)
	return resourceOpenAIFileRead(ctx, d, m)
}

func resourceOpenAIFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	fileObj, err := client.Files.Get(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("bytes", fileObj.Bytes)
	d.Set("created_at", fileObj.CreatedAt)
	d.Set("filename", fileObj.Filename)
	d.Set("purpose", string(fileObj.Purpose))

	return nil
}

func resourceOpenAIFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	_, err := client.Files.Delete(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
