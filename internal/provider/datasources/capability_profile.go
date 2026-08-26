package datasources

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type CapabilityProfileDataSource struct{}

type CapabilityProfileDataSourceModel struct {
	Capability     types.String `tfsdk:"capability"`
	Profile        types.String `tfsdk:"profile"`
	CatalogVersion types.String `tfsdk:"catalog_version"`
	Version        types.String `tfsdk:"version"`
	Service        types.String `tfsdk:"service"`
	Configuration  types.String `tfsdk:"configuration"`
}

func (d CapabilityProfileDataSource) Metadata(
	ctx context.Context,
	req datasource.MetadataRequest,
	resp *datasource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_capability_profile"
}

func (d CapabilityProfileDataSource) Schema(
	ctx context.Context,
	req datasource.SchemaRequest,
	resp *datasource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"capability": schema.StringAttribute{
				Required: true,
			},
			"profile": schema.StringAttribute{
				Required: true,
			},
			"catalog_version": schema.StringAttribute{
				Computed: true,
			},
			"version": schema.StringAttribute{
				Computed: true,
			},
			"service": schema.StringAttribute{
				Computed: true,
			},
			"configuration": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d CapabilityProfileDataSource) Read(
	ctx context.Context,
	req datasource.ReadRequest,
	resp *datasource.ReadResponse,
) {
	var data CapabilityProfileDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	data.CatalogVersion = types.StringValue("1.1")

	switch data.Capability.ValueString() {
	case "object-storage":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("s3")

		switch data.Profile.ValueString() {
		case "temporary":
			data.Configuration = types.StringValue(
				`{"retention":{"strategy":"temporary"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`,
			)

		case "standard":
			data.Configuration = types.StringValue(
				`{"retention":{"strategy":"retained"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`,
			)

		case "archive":
			data.Configuration = types.StringValue(
				`{"retention":{"strategy":"archive"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`,
			)

		case "compliance_archive":
			data.Configuration = types.StringValue(
				`{"retention":{"strategy":"indefinite"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}}`,
			)

		default:
			resp.Diagnostics.AddError(
				"Capability profile not found",
				fmt.Sprintf(
					"The profile %q is not available for capability %q.",
					data.Profile.ValueString(),
					data.Capability.ValueString(),
				),
			)
			return
		}

	case "cryptographic-key-management":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("kms")

		switch data.Profile.ValueString() {
		case "data_encryption":
			data.Configuration = types.StringValue(
				`{"cryptography":{"type":"symmetric","usage":"encrypt_decrypt","specification":"symmetric_default"},"lifecycle":{"rotation":{"enabled":true,"period_in_days":365}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`,
			)

		case "hmac":
			data.Configuration = types.StringValue(
				`{"cryptography":{"type":"symmetric","usage":"generate_verify_mac","specification":"hmac_256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`,
			)

		case "signing":
			data.Configuration = types.StringValue(
				`{"cryptography":{"type":"asymmetric","usage":"sign_verify","specification":"rsa_3072"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`,
			)

		case "key_agreement":
			data.Configuration = types.StringValue(
				`{"cryptography":{"type":"asymmetric","usage":"key_agreement","specification":"ecc_nist_p256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`,
			)

		default:
			resp.Diagnostics.AddError(
				"Capability profile not found",
				fmt.Sprintf(
					"The profile %q is not available for capability %q.",
					data.Profile.ValueString(),
					data.Capability.ValueString(),
				),
			)
			return
		}

	case "audit-logging":
		data.Version = types.StringValue("1.0")
		data.Service = types.StringValue("audit_logging")

		switch data.Profile.ValueString() {
		case "standard":
			data.Configuration = types.StringValue(
				`{"collection":{"scope":"workload","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`,
			)

		case "centralized":
			data.Configuration = types.StringValue(
				`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`,
			)

		case "security":
			data.Configuration = types.StringValue(
				`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"extended"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"GOVERNANCE","retention_days":365}},"access":{"public_access":false,"controlled":true}}`,
			)

		case "compliance":
			data.Configuration = types.StringValue(
				`{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"indefinite"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}},"access":{"public_access":false,"controlled":true}}`,
			)

		default:
			resp.Diagnostics.AddError(
				"Capability profile not found",
				fmt.Sprintf(
					"The profile %q is not available for capability %q.",
					data.Profile.ValueString(),
					data.Capability.ValueString(),
				),
			)
			return
		}

	default:
		resp.Diagnostics.AddError(
			"Capability not found",
			fmt.Sprintf(
				"The capability %q is not available.",
				data.Capability.ValueString(),
			),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func NewCapabilityProfileDataSource() datasource.DataSource {
	return CapabilityProfileDataSource{}
}

func resolveCapabilityProfile(
	capability string,
	profile string,
) (version string, service string, configuration string, err error) {
	switch capability {
	case "object-storage":
		version = "1.0"
		service = "s3"

		switch profile {
		case "temporary":
			configuration = `{"retention":{"strategy":"temporary"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`
		case "standard":
			configuration = `{"retention":{"strategy":"retained"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`
		case "archive":
			configuration = `{"retention":{"strategy":"archive"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":false}}`
		case "compliance_archive":
			configuration = `{"retention":{"strategy":"indefinite"},"encryption":{"mode":"sse_s3"},"immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}`
		default:
			return "", "", "", fmt.Errorf(
				"the profile %q is not available for capability %q",
				profile,
				capability,
			)
		}

	case "cryptographic-key-management":
		version = "1.0"
		service = "kms"

		switch profile {
		case "data_encryption":
			configuration = `{"cryptography":{"type":"symmetric","usage":"encrypt_decrypt","specification":"symmetric_default"},"lifecycle":{"rotation":{"enabled":true,"period_in_days":365}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`
		case "hmac":
			configuration = `{"cryptography":{"type":"symmetric","usage":"generate_verify_mac","specification":"hmac_256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional","multi_region_primary","multi_region_replica"]}}`
		case "signing":
			configuration = `{"cryptography":{"type":"asymmetric","usage":"sign_verify","specification":"rsa_3072"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`
		case "key_agreement":
			configuration = `{"cryptography":{"type":"asymmetric","usage":"key_agreement","specification":"ecc_nist_p256"},"lifecycle":{"rotation":{"enabled":false}},"deployment":{"allowed_scopes":["regional"]}}`
		default:
			return "", "", "", fmt.Errorf(
				"the profile %q is not available for capability %q",
				profile,
				capability,
			)
		}

	case "audit-logging":
		version = "1.0"
		service = "audit_logging"

		switch profile {
		case "standard":
			configuration = `{"collection":{"scope":"workload","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`
		case "centralized":
			configuration = `{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"retained"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":false}},"access":{"public_access":false,"controlled":true}}`
		case "security":
			configuration = `{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"extended"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"GOVERNANCE","retention_days":365}},"access":{"public_access":false,"controlled":true}}`
		case "compliance":
			configuration = `{"collection":{"scope":"organization","multi_region":true},"retention":{"strategy":"indefinite"},"protection":{"integrity_validation":true,"encryption":"customer_managed","immutability":{"enabled":true,"mode":"COMPLIANCE","retention_days":2555}},"access":{"public_access":false,"controlled":true}}`
		default:
			return "", "", "", fmt.Errorf(
				"the profile %q is not available for capability %q",
				profile,
				capability,
			)
		}

	default:
		return "", "", "", fmt.Errorf(
			"the capability %q is not available",
			capability,
		)
	}

	return version, service, configuration, nil
}
