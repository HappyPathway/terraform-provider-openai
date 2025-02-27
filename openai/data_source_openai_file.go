package openai

import (
	"context"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func dataSourceOpenAIFile() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceOpenAIFileRead,
		Schema: map[string]*schema.Schema{
			"file_id": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The ID of the file to retrieve",
			},
			"bytes": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "Size of the file in bytes",
			},
			"filename": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "Name of the file",
			},
			"purpose": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The intended purpose of the file",
			},
			"created_at": {
				Type:        schema.TypeInt,
				Computed:    true,
				Description: "The Unix timestamp for when the file was created",
			},
		},
	}
}

func dataSourceOpenAIFileRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	client := m.(*Client)

	fileID := d.Get("file_id").(string)

	file, err := client.GetFile(ctx, fileID)
	if err != nil {
		return diag.FromErr(err)
	}

	d.SetId(file.ID)
	d.Set("bytes", file.Bytes)
	d.Set("filename", file.Filename)
	d.Set("purpose", file.Purpose)
	d.Set("created_at", file.CreatedAt)

	return nil
}
