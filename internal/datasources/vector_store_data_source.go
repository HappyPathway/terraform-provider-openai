package datasources

import (
	"context"
	"fmt"

	"github.com/darnold/terraform-provider-openai/internal/client"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure provider defined types fully satisfy framework interfaces
var _ datasource.DataSource = &VectorStoreDataSource{}

func NewVectorStoreDataSource() datasource.DataSource {
	return &VectorStoreDataSource{}
}

// VectorStoreDataSource defines the data source implementation.
type VectorStoreDataSource struct {
	client *client.Client
}

// VectorStoreDataSourceModel describes the data source data model.
type VectorStoreDataSourceModel struct {
	ID         types.String `tfsdk:"id"`
	Name       types.String `tfsdk:"name"`
	CreatedAt  types.Int64  `tfsdk:"created_at"`
	Status     types.String `tfsdk:"status"`
	FileCounts struct {
		InProgress types.Int64 `tfsdk:"in_progress"`
		Completed  types.Int64 `tfsdk:"completed"`
		Failed     types.Int64 `tfsdk:"failed"`
		Cancelled  types.Int64 `tfsdk:"cancelled"`
		Total      types.Int64 `tfsdk:"total"`
	} `tfsdk:"file_counts"`
	UsageBytes   types.Int64 `tfsdk:"usage_bytes"`
	ExpiresAt    types.Int64 `tfsdk:"expires_at"`
	ExpiresAfter *struct {
		Days   types.Int64  `tfsdk:"days"`
		Anchor types.String `tfsdk:"anchor"`
	} `tfsdk:"expires_after"`
	Metadata types.Map `tfsdk:"metadata"`
}

func (d *VectorStoreDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vector_store"
}

func (d *VectorStoreDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		MarkdownDescription: "Use this data source to get information about an existing OpenAI vector store.",

		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				MarkdownDescription: "The ID of the vector store to retrieve.",
				Required:            true,
			},
			"name": schema.StringAttribute{
				MarkdownDescription: "The name of the vector store.",
				Computed:            true,
			},
			"created_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the vector store was created.",
				Computed:            true,
			},
			"status": schema.StringAttribute{
				MarkdownDescription: "The current status of the vector store.",
				Computed:            true,
			},
			"file_counts": schema.SingleNestedAttribute{
				MarkdownDescription: "Statistics about files in the vector store.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"in_progress": schema.Int64Attribute{
						MarkdownDescription: "Number of files currently being processed.",
						Computed:            true,
					},
					"completed": schema.Int64Attribute{
						MarkdownDescription: "Number of successfully processed files.",
						Computed:            true,
					},
					"failed": schema.Int64Attribute{
						MarkdownDescription: "Number of files that failed processing.",
						Computed:            true,
					},
					"cancelled": schema.Int64Attribute{
						MarkdownDescription: "Number of cancelled file operations.",
						Computed:            true,
					},
					"total": schema.Int64Attribute{
						MarkdownDescription: "Total number of files.",
						Computed:            true,
					},
				},
			},
			"usage_bytes": schema.Int64Attribute{
				MarkdownDescription: "The total size of the vector store in bytes.",
				Computed:            true,
			},
			"expires_at": schema.Int64Attribute{
				MarkdownDescription: "The timestamp when the vector store will expire.",
				Computed:            true,
			},
			"expires_after": schema.SingleNestedAttribute{
				MarkdownDescription: "Configuration for vector store expiration.",
				Computed:            true,
				Attributes: map[string]schema.Attribute{
					"days": schema.Int64Attribute{
						MarkdownDescription: "Number of days after which the vector store expires.",
						Computed:            true,
					},
					"anchor": schema.StringAttribute{
						MarkdownDescription: "Reference time for expiration calculation.",
						Computed:            true,
					},
				},
			},
			"metadata": schema.MapAttribute{
				MarkdownDescription: "Metadata key-value pairs for the vector store.",
				Computed:            true,
				ElementType:         types.StringType,
			},
		},
	}
}

func (d *VectorStoreDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	client, ok := req.ProviderData.(*client.Client)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Data Source Configure Type",
			fmt.Sprintf("Expected *client.Client, got: %T", req.ProviderData),
		)
		return
	}

	d.client = client
}

func (d *VectorStoreDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var data VectorStoreDataSourceModel
	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	result, err := d.client.OpenAI.RetrieveVectorStore(ctx, data.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Reading Vector Store",
			fmt.Sprintf("Unable to read vector store: %s", d.client.HandleError(err)),
		)
		return
	}

	// Map response to model
	data.Name = types.StringValue(result.Name)
	data.CreatedAt = types.Int64Value(result.CreatedAt)
	data.Status = types.StringValue(result.Status)
	data.FileCounts.InProgress = types.Int64Value(int64(result.FileCounts.InProgress))
	data.FileCounts.Completed = types.Int64Value(int64(result.FileCounts.Completed))
	data.FileCounts.Failed = types.Int64Value(int64(result.FileCounts.Failed))
	data.FileCounts.Cancelled = types.Int64Value(int64(result.FileCounts.Cancelled))
	data.FileCounts.Total = types.Int64Value(int64(result.FileCounts.Total))
	data.UsageBytes = types.Int64Value(int64(result.UsageBytes))

	if result.ExpiresAt != nil {
		data.ExpiresAt = types.Int64Value(int64(*result.ExpiresAt))
	}

	if result.ExpiresAfter != nil {
		data.ExpiresAfter = &struct {
			Days   types.Int64  `tfsdk:"days"`
			Anchor types.String `tfsdk:"anchor"`
		}{
			Days:   types.Int64Value(int64(result.ExpiresAfter.Days)),
			Anchor: types.StringValue(result.ExpiresAfter.Anchor),
		}
	}

	// Convert metadata
	metadataMap := make(map[string]attr.Value)
	for k, v := range result.Metadata {
		metadataMap[k] = types.StringValue(fmt.Sprintf("%v", v))
	}
	metadata, diags := types.MapValue(types.StringType, metadataMap)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	data.Metadata = metadata

	// Save into state
	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}
