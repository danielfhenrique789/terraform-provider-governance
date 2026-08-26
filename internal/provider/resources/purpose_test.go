package resources

import (
	"context"
	"testing"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/definitions"
	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
	frameworkresource "github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	
)

func TestPurposeResourceConfigure(t *testing.T) {
	repository := &definitions.MockRepository{
		Catalog: purpose.PurposeCatalog{
			Version: "1.0",
			Purposes: map[string]purpose.Purpose{
				"centralized_audit_logging": {
					Name: "centralized_audit_logging",
					Capabilities: map[string]string{
						"audit": "standard",
					},
				},
			},
		},
	}

	purposeResource := &PurposeResource{}

	var response frameworkresource.ConfigureResponse

	purposeResource.Configure(
		context.Background(),
		frameworkresource.ConfigureRequest{
			ProviderData: repository,
		},
		&response,
	)

	if response.Diagnostics.HasError() {
		t.Fatalf(
			"unexpected diagnostics: %v",
			response.Diagnostics,
		)
	}

	if purposeResource.repository != repository {
		t.Fatal("repository was not configured")
	}
}

func TestPurposeResourceConfigureInvalidProviderData(t *testing.T) {
	purposeResource := &PurposeResource{}

	var response frameworkresource.ConfigureResponse

	purposeResource.Configure(
		context.Background(),
		frameworkresource.ConfigureRequest{
			ProviderData: "invalid",
		},
		&response,
	)

	if !response.Diagnostics.HasError() {
		t.Fatal("expected diagnostics")
	}
}

func TestPurposeResourceConfigureNilProviderData(t *testing.T) {
	purposeResource := &PurposeResource{}

	var response frameworkresource.ConfigureResponse

	purposeResource.Configure(
		context.Background(),
		frameworkresource.ConfigureRequest{},
		&response,
	)

	if response.Diagnostics.HasError() {
		t.Fatalf(
			"unexpected diagnostics: %v",
			response.Diagnostics,
		)
	}

	if purposeResource.repository != nil {
		t.Fatal("repository should remain nil")
	}
}

func TestPurposeResourceModel(t *testing.T) {
	data := PurposeResourceModel{
		Name: types.StringValue("centralized_audit_logging"),
	}

	if data.Name.ValueString() != "centralized_audit_logging" {
		t.Fatalf(
			"unexpected purpose name: %q",
			data.Name.ValueString(),
		)
	}
}

func TestPurposeResourceSchema(t *testing.T) {
	purposeResource := &PurposeResource{}

	var response frameworkresource.SchemaResponse

	purposeResource.Schema(
		context.Background(),
		frameworkresource.SchemaRequest{},
		&response,
	)

	if response.Diagnostics.HasError() {
		t.Fatalf(
			"unexpected diagnostics: %v",
			response.Diagnostics,
		)
	}

	name, exists := response.Schema.Attributes["name"]
	if !exists {
		t.Fatal("name attribute is missing")
	}

	if !name.IsRequired() {
		t.Fatal("name attribute should be required")
	}

	capabilities, exists := response.Schema.Attributes["capabilities"]
	if !exists {
		t.Fatal("capabilities attribute is missing")
	}

	if !capabilities.IsComputed() {
		t.Fatal("capabilities attribute should be computed")
	}
}