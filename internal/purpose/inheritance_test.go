package purpose

import "testing"

func TestPurposeCatalogValidateInheritance(t *testing.T) {
	parent := "basic_data_storage"

	tests := []struct {
		name    string
		catalog PurposeCatalog
		wantErr bool
	}{
		{
			name: "valid inheritance",
			catalog: PurposeCatalog{
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
							"retention": "compliance",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "unknown parent",
			catalog: PurposeCatalog{
				Version: "1.0",
				Purposes: map[string]Purpose{
					"compliant_data_storage": {
						Name:        "compliant_data_storage",
						Description: "Compliant enterprise data storage.",
						Extends:     stringPtr("unknown"),
						Capabilities: map[string]string{
							"retention": "compliance",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "circular inheritance",
			catalog: PurposeCatalog{
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
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.catalog.ValidateInheritance()

			if (err != nil) != tt.wantErr {
				t.Fatalf(
					"ValidateInheritance() error = %v, wantErr = %v",
					err,
					tt.wantErr,
				)
			}
		})
	}
}

func stringPtr(value string) *string {
	return &value
}