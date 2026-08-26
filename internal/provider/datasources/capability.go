package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CapabilityDataSource struct{}

type CapabilityDataSourceModel struct {
	Name           types.String `tfsdk:"name"`
	Version        types.String `tfsdk:"version"`
	Service        types.String `tfsdk:"service"`
	Profiles       types.Map    `tfsdk:"profiles"`
	CatalogVersion types.String `tfsdk:"catalog_version"`
}

func (d CapabilityDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_capability"
}

func (d CapabilityDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"service": schema.StringAttribute{
				Computed: true,
			},
			"profiles": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"catalog_version": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d CapabilityDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data CapabilityDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.CatalogVersion = types.StringValue("1.0.1")
	switch data.Name.ValueString() {
	case "object_storage":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("s3")

		data.Profiles = types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"temporary": types.StringValue(
					`{"retention":{"strategy":"temporary"},"encryption":{"mode":"sse_s3"}}`,
				),
				"standard": types.StringValue(
					`{"retention":{"strategy":"retained"},"encryption":{"mode":"sse_s3"}}`,
				),
				"archive": types.StringValue(
					`{"retention":{"strategy":"archive"},"encryption":{"mode":"sse_s3"}}`,
				),
				"compliance_archive": types.StringValue(
					`{"retention":{"strategy":"indefinite"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}}`,
				),
			},
		)

	case "cryptographic-key-management":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("kms")

		data.Profiles = types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"data_encryption": types.StringValue(
					`{"cryptography":{"type":"symmetric","usage":"encrypt_decrypt","specification":"symmetric_default"},"lifecycle":{"rotation":{"enabled":true,"period_in_days":365}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`,
				),
				"hmac": types.StringValue(
					`{"cryptography":{"type":"symmetric","usage":"generate_verify_mac","specification":"hmac_256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`,
				),
				"signing": types.StringValue(
					`{"cryptography":{"type":"asymmetric","usage":"sign_verify","specification":"rsa_3072"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`,
				),
				"key_agreement": types.StringValue(
					`{"cryptography":{"type":"asymmetric","usage":"key_agreement","specification":"ecc_nist_p256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`,
				),
			},
		)

	case "audit_logging":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("audit_logging")

		data.Profiles = types.MapValueMust(
			types.StringType,
			map[string]attr.Value{
				"standard": types.StringValue(
					`{"collection":{"scope":"workload","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`,
				),
				"centralized": types.StringValue(
					`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`,
				),
				"security": types.StringValue(
					`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"extended"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"GOVERNANCE","retention_days":365}},"access":{"public_access":false,"controlled":true}}`,
				),
				"compliance": types.StringValue(
					`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"indefinite"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}},"access":{"public_access":false,"controlled":true}}`,
				),
			},
		)

	default:
		resp.Diagnostics.AddError(
			"Capability not found",
			fmt.Sprintf(
				"The capability %q is not available.",
				data.Name.ValueString(),
			),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewCapabilityDataSource() datasource.DataSource {
	return CapabilityDataSource{}
}
