package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PurposeDataSource struct{}

type PurposeDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Description    types.String `tfsdk:"description"`
	Capabilities   types.Map    `tfsdk:"capabilities"`
	CatalogVersion types.String `tfsdk:"catalog_version"`
}

func (d PurposeDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_purpose"
}

func (d PurposeDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"capabilities": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"catalog_version": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d PurposeDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data PurposeDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	resolved, err := resolvePurpose(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Purpose not found",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(
		resp.State.Set(ctx, &resolved)...,
	)
}

func NewPurposeDataSource() datasource.DataSource {
	return PurposeDataSource{}
}

func resolvePurpose(name string) (PurposeDataSourceModel, error) {
	data := PurposeDataSourceModel{
		CatalogVersion: types.StringValue("1.0.1"),
	}

	switch name {
	case "centralized_audit_logging":
		data.Name = types.StringValue(name)
		data.Description = types.StringValue(
			"Centralized, protected retention of organization-wide audit and configuration logs for security operations, governance, and compliance.",
		)

		data.Capabilities = types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"object_storage":                types.StringValue("compliance_archive"),
				"cryptographic-key-management": types.StringValue("data_encryption"),
				"audit_logging":                 types.StringValue("centralized"),
			},
		)

	case "simple_data_usage":
		data.Name = types.StringValue(name)
		data.Description = types.StringValue(
			"This should declare only an object storage without any extra detail.",
		)

		data.Capabilities = types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"object_storage": types.StringValue("simple_use"),
			},
		)

	default:
		return PurposeDataSourceModel{}, fmt.Errorf(
			"purpose %q is not available",
			name,
		)
	}

	return data, nil
}