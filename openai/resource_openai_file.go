package openai

import (
	"context"
	"io/ioutil"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
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
	client := m.(*Client)

	filePath := d.Get("file").(string)
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		return diag.FromErr(err)
	}

	uploadReq := &FileUploadRequest{
		File:    fileContent,
		Purpose: d.Get("purpose").(string),
	}

	file, err := client.UploadFile(ctx, uploadReq)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(file.ID)

	return resourceOpenAIFileRead(ctx, d, m)
}

func resourceOpenAIFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	file, err := client.GetFile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.Set("bytes", file.Bytes)
	d.Set("created_at", file.CreatedAt)
	d.Set("filename", file.Filename)
	d.Set("purpose", file.Purpose)

	return nil
}

func resourceOpenAIFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	err := client.DeleteFile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")
	return nil
}
