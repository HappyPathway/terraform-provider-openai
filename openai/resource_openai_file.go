package openai

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/HappyPathway/terraform-provider-openai/openai/testutil"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func resourceOpenAIFile() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourceOpenAIFileCreate,
		ReadContext:   resourceOpenAIFileRead,
		DeleteContext: resourceOpenAIFileDelete,
		Schema: map[string]*schema.Schema{
			"content": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The base64-encoded content of the file. Use the filebase64() function to read a file from disk.",
			},
			"filename": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The name of the file to be uploaded.",
			},
			"purpose": {
				Type:        schema.TypeString,
				Required:    true,
				ForceNew:    true,
				Description: "The intended purpose of the uploaded file. Valid values are 'fine-tune', 'assistants', or 'fine-tune-results'.",
			},
			"bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The size of the file in bytes.",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp (in seconds) for when the file was created.",
			},
			"object": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The object type, which is always 'file'.",
			},
		},
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
	}
}

func resourceOpenAIFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(testutil.ClientInterface)

	content := d.Get("content").(string)
	fileBytes, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error decoding file content: %v", err))
	}

	req := &testutil.FileUploadRequest{
		File:     fileBytes,
		Purpose:  d.Get("purpose").(string),
		Filename: d.Get("filename").(string),
	}

	file, err := client.UploadFile(ctx, req)
	if err != nil {
		return diag.FromErr(fmt.Errorf("error uploading file: %v", err))
	}

	d.SetId(file.ID)
	d.Set("bytes", file.Bytes)
	d.Set("created_at", file.CreatedAt)
	d.Set("object", file.Object)

	return nil
}

func resourceOpenAIFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	file, err := client.GetFile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error retrieving file: %v", err))
	}

	d.Set("filename", file.Filename)
	d.Set("purpose", file.Purpose)
	d.Set("bytes", file.Bytes)
	d.Set("created_at", file.CreatedAt)
	d.Set("object", file.Object)

	return nil
}

func resourceOpenAIFileDelete(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	err := client.DeleteFile(ctx, d.Id())
	if err != nil {
		return diag.FromErr(fmt.Errorf("error deleting file: %v", err))
	}

	d.SetId("")
	return nil
}
