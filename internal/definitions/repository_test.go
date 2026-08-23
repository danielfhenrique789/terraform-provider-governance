package definitions_test

import (
	"context"
	"testing"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/definitions"
	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
)

func TestMockRepositoryGetPurposeCatalog(t *testing.T) {
	repository := &definitions.MockRepository{
		Catalog: purpose.PurposeCatalog{
			Version: "1.0",
			Purposes: map[string]purpose.Purpose{
				"customer_document_storage": {
					Name:        "customer_document_storage",
					Description: "Persistent storage for customer documents.",
					Capabilities: map[string]string{
						"object_storage": "standard",
					},
				},
			},
		},
	}

	catalog, err := repository.GetPurposeCatalog(context.Background())
	if err != nil {
		t.Fatalf("GetPurposeCatalog() error = %v", err)
	}

	if catalog.Version != "1.0" {
		t.Fatalf("expected version 1.0, got %q", catalog.Version)
	}

	p, exists := catalog.Purposes["customer_document_storage"]
	if !exists {
		t.Fatal("expected customer_document_storage purpose")
	}

	if p.Capabilities["object_storage"] != "standard" {
		t.Fatalf("unexpected object_storage profile")
	}
}