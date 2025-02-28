package openai

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

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
			"file_path": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Path to the file to upload",
				ExactlyOneOf: []string{"file_path", "content"},
			},
			"content": {
				Type:         schema.TypeString,
				Optional:     true,
				ForceNew:     true,
				Description:  "Direct content to upload as a file",
				ExactlyOneOf: []string{"file_path", "content"},
			},
			"filename": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				ForceNew:    true,
				Description: "Name of the file. Required when using content field, computed when using file_path.",
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
		},
	}
}

func resourceOpenAIFileCreate(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	config := m.(*Config)
	client := config.Client

	var reader io.Reader
	var filename string

	if filePath, ok := d.GetOk("file_path"); ok {
		file, err := os.Open(filePath.(string))
		if err != nil {
			return diag.FromErr(err)
		}
		defer file.Close()
		reader = file
		filename = filePath.(string)
	} else if content, ok := d.GetOk("content"); ok {
		reader = strings.NewReader(content.(string))
		filename = d.Get("filename").(string)
		if filename == "" {
			return diag.FromErr(fmt.Errorf("filename is required when using content field"))
		}
	}

	fileObj, err := client.Files.New(ctx, openaiapi.FileNewParams{
		File:    openaiapi.FileParam(reader, filename, "application/octet-stream"),
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
