package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestPurposeDataSource_Metadata(t *testing.T) {
	d := PurposeDataSource{}

	var resp datasource.MetadataResponse

	d.Metadata(
		context.Background(),
		datasource.MetadataRequest{
			ProviderTypeName: "governance",
		},
		&resp,
	)

	if resp.TypeName != "governance_purpose" {
		t.Fatalf(
			"expected type name %q, got %q",
			"governance_purpose",
			resp.TypeName,
		)
	}
}

func TestPurposeDataSource_Schema(t *testing.T) {
	d := PurposeDataSource{}

	var resp datasource.SchemaResponse

	d.Schema(
		context.Background(),
		datasource.SchemaRequest{},
		&resp,
	)

	if resp.Diagnostics.HasError() {
		t.Fatalf(
			"unexpected schema diagnostics: %v",
			resp.Diagnostics,
		)
	}

	for _, name := range []string{
		"name",
		"description",
		"capabilities",
		"catalog_version",
	} {
		if _, ok := resp.Schema.Attributes[name]; !ok {
			t.Errorf("expected schema attribute %q", name)
		}
	}
}

func TestResolvePurpose_CentralizedAuditLogging(t *testing.T) {
	data, err := resolvePurpose("centralized_audit_logging")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Name.ValueString() != "centralized_audit_logging" {
		t.Fatalf(
			"expected purpose name %q, got %q",
			"centralized_audit_logging",
			data.Name.ValueString(),
		)
	}

	if data.CatalogVersion.ValueString() != "1.0.1" {
		t.Fatalf(
			"expected catalog version %q, got %q",
			"1.0.1",
			data.CatalogVersion.ValueString(),
		)
	}

	expected := map[string]string{
		"object_storage":                "compliance_archive",
		"cryptographic-key-management": "data_encryption",
		"audit_logging":                 "centralized",
	}

	for capability, profile := range expected {
		value, ok := data.Capabilities.Elements()[capability]

		if !ok {
			t.Errorf(
				"expected capability %q",
				capability,
			)
			continue
		}

		if value.String() != `"`+profile+`"` {
			t.Errorf(
				"expected capability %q to use profile %q, got %s",
				capability,
				profile,
				value.String(),
			)
		}
	}
}

func TestResolvePurpose_SimpleDataUsage(t *testing.T) {
	data, err := resolvePurpose("simple_data_usage")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if data.Name.ValueString() != "simple_data_usage" {
		t.Fatalf(
			"expected purpose name %q, got %q",
			"simple_data_usage",
			data.Name.ValueString(),
		)
	}

	if len(data.Capabilities.Elements()) != 1 {
		t.Fatalf(
			"expected 1 capability, got %d",
			len(data.Capabilities.Elements()),
		)
	}

	value, ok := data.Capabilities.Elements()["object_storage"]

	if !ok {
		t.Fatal("expected object_storage capability")
	}

	if value.String() != `"simple_use"` {
		t.Fatalf(
			"expected object_storage profile %q, got %s",
			"simple_use",
			value.String(),
		)
	}
}

func TestResolvePurpose_UnknownPurpose(t *testing.T) {
	_, err := resolvePurpose("does_not_exist")

	if err == nil {
		t.Fatal("expected error for unknown purpose")
	}
}