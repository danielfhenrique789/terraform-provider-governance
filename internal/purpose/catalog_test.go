package purpose

import "testing"

func TestPurposeCatalogValidate(t *testing.T) {
	tests := []struct {
		name    string
		catalog PurposeCatalog
		wantErr bool
	}{
		{
			name: "valid catalog",
			catalog: PurposeCatalog{
				Version: "1.0",
				Purposes: map[string]Purpose{
					"customer_document_storage": {
						Name:        "customer_document_storage",
						Description: "Persistent storage for customer documents.",
						Capabilities: map[string]string{
							"object_storage": "standard",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "missing version",
			catalog: PurposeCatalog{
				Purposes: map[string]Purpose{
					"customer_document_storage": {
						Name:        "customer_document_storage",
						Description: "Persistent storage for customer documents.",
						Capabilities: map[string]string{
							"object_storage": "standard",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "no purposes",
			catalog: PurposeCatalog{
				Version:  "1.0",
				Purposes: map[string]Purpose{},
			},
			wantErr: true,
		},
		{
			name: "purpose name does not match catalog key",
			catalog: PurposeCatalog{
				Version: "1.0",
				Purposes: map[string]Purpose{
					"customer_document_storage": {
						Name:        "different_name",
						Description: "Persistent storage for customer documents.",
						Capabilities: map[string]string{
							"object_storage": "standard",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.catalog.Validate()

			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestPurposeCatalogValidateRejectsCircularInheritance(t *testing.T) {
	parentA := "purpose_b"
	parentB := "purpose_a"

	catalog := PurposeCatalog{
		Version: "1.0",
		Purposes: map[string]Purpose{
			"purpose_a": {
				Name:        "purpose_a",
				Extends:     &parentA,
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			"purpose_b": {
				Name:        "purpose_b",
				Extends:     &parentB,
				Capabilities: map[string]string{
					"retention": "compliance",
				},
			},
		},
	}

	if err := catalog.Validate(); err == nil {
		t.Fatal("expected circular inheritance validation error")
	}
}