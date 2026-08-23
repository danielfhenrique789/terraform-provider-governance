package purpose

import "testing"

func TestPurposeCatalogResolve(t *testing.T) {
	parent := "basic_data_storage"

	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"basic_data_storage": {
				Name:        "basic_data_storage",
				Description: "Basic enterprise data storage.",
				Capabilities: map[string]string{
					"object_storage": "standard",
					"key_management": "data_encryption",
				},
			},
			"compliant_data_storage": {
				Name:        "compliant_data_storage",
				Description: "Compliant enterprise data storage.",
				Extends:     &parent,
				Capabilities: map[string]string{
					"retention":    "compliance",
					"immutability": "compliance",
				},
			},
		},
	}

	resolved, err := catalog.Resolve("compliant_data_storage")
	if err != nil {
		t.Fatalf("Resolve() error = %v", err)
	}

	expected := map[string]string{
		"object_storage": "standard",
		"key_management": "data_encryption",
		"retention":      "compliance",
		"immutability":   "compliance",
	}

	if len(resolved.Capabilities) != len(expected) {
		t.Fatalf(
			"expected %d capabilities, got %d",
			len(expected),
			len(resolved.Capabilities),
		)
	}

	for capability, expectedProfile := range expected {
		actualProfile, exists := resolved.Capabilities[capability]

		if !exists {
			t.Fatalf("missing capability %q", capability)
		}

		if actualProfile != expectedProfile {
			t.Fatalf(
				"capability %q: expected profile %q, got %q",
				capability,
				expectedProfile,
				actualProfile,
			)
		}
	}
}

func TestPurposeCatalogResolveRejectsCapabilityConflict(t *testing.T) {
	parent := "basic_data_storage"

	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"basic_data_storage": {
				Name:        "basic_data_storage",
				Description: "Basic enterprise data storage.",
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			"compliant_data_storage": {
				Name:        "compliant_data_storage",
				Description: "Compliant enterprise data storage.",
				Extends:     &parent,
				Capabilities: map[string]string{
					"object_storage": "archive",
				},
			},
		},
	}

	_, err := catalog.Resolve("compliant_data_storage")

	if err == nil {
		t.Fatal("expected capability conflict error")
	}
}

func TestPurposeCatalogResolveRejectsUnknownPurpose(t *testing.T) {
	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"customer_data": {
				Name:        "customer_data",
				Description: "Customer data.",
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
		},
	}

	_, err := catalog.Resolve("does_not_exist")

	if err == nil {
		t.Fatal("expected unknown purpose error")
	}
}

func TestPurposeCatalogResolveRejectsUnknownParent(t *testing.T) {
	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"customer_data": {
				Name:        "customer_data",
				Description: "Customer data.",
				Extends:     stringPtr("does_not_exist"),
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
		},
	}

	_, err := catalog.Resolve("customer_data")

	if err == nil {
		t.Fatal("expected unknown parent error")
	}
}

func TestPurposeCatalogResolveRejectsCircularInheritance(t *testing.T) {
	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"purpose_a": {
				Name:        "purpose_a",
				Description: "Purpose A.",
				Extends:     stringPtr("purpose_b"),
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			"purpose_b": {
				Name:        "purpose_b",
				Description: "Purpose B.",
				Extends:     stringPtr("purpose_a"),
				Capabilities: map[string]string{
					"retention": "compliance",
				},
			},
		},
	}

	_, err := catalog.Resolve("purpose_a")

	if err == nil {
		t.Fatal("expected circular inheritance error")
	}
}