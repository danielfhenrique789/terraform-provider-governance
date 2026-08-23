package definitions

import (
	"testing"
)

func TestLoadPurposeCatalog(t *testing.T) {
	data := []byte(`
version: "1.0"

purposes:
  customer_document_storage:
    description: >
      Persistent storage for customer documents.

    capabilities:
      object_storage: standard
      key_management: data_encryption
`)

	catalog, err := LoadPurposeCatalog(data)
	if err != nil {
		t.Fatalf("LoadPurposeCatalog() error = %v", err)
	}

	if catalog.Version != "1.0" {
		t.Fatalf("expected version 1.0, got %q", catalog.Version)
	}

	p, exists := catalog.Purposes["customer_document_storage"]
	if !exists {
		t.Fatal("expected customer_document_storage purpose")
	}

	if p.Name != "customer_document_storage" {
		t.Fatalf("unexpected purpose name: %q", p.Name)
	}

	if p.Description == "" {
		t.Fatal("expected purpose description")
	}

	if p.Capabilities["object_storage"] != "standard" {
		t.Fatalf("unexpected object_storage profile")
	}

	if p.Capabilities["key_management"] != "data_encryption" {
		t.Fatalf("unexpected key_management profile")
	}
}

func TestLoadPurposeCatalogRejectsPurposeWithoutCapabilities(t *testing.T) {
	data := []byte(`
version: "1.0"

purposes:
  invalid_purpose:
    description: >
      This purpose has no capabilities.
`)

	_, err := LoadPurposeCatalog(data)

	if err == nil {
		t.Fatal("expected validation error")
	}
}

func TestLoadPurposeCatalogPreservesPurposeInheritance(t *testing.T) {
	data := []byte(`
version: "1.0"

purposes:
  basic_data_storage:
    description: >
      Basic enterprise data storage.

    capabilities:
      object_storage: standard

  compliant_data_storage:
    description: >
      Compliant enterprise data storage.

    extends: basic_data_storage

    capabilities:
      retention: compliance
`)

	catalog, err := LoadPurposeCatalog(data)
	if err != nil {
		t.Fatalf("LoadPurposeCatalog() error = %v", err)
	}

	p, exists := catalog.Purposes["compliant_data_storage"]
	if !exists {
		t.Fatal("expected compliant_data_storage purpose")
	}

	if p.Extends == nil {
		t.Fatal("expected extends to be set")
	}

	if *p.Extends != "basic_data_storage" {
		t.Fatalf(
			"expected parent basic_data_storage, got %q",
			*p.Extends,
		)
	}
}
