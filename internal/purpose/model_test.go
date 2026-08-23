package purpose

import "testing"

func TestPurposeValidate(t *testing.T) {
	tests := []struct {
		name    string
		purpose Purpose
		wantErr bool
	}{
		{
			name: "valid purpose",
			purpose: Purpose{
				Name:        "customer_document_storage",
				Description: "Persistent storage for customer documents.",
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			wantErr: false,
		},
		{
			name: "missing name",
			purpose: Purpose{
				Description: "Persistent storage for customer documents.",
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			wantErr: true,
		},
		{
			name: "missing description",
			purpose: Purpose{
				Name: "customer_document_storage",
				Capabilities: map[string]string{
					"object_storage": "standard",
				},
			},
			wantErr: true,
		},
		{
			name: "no capabilities",
			purpose: Purpose{
				Name:         "customer_document_storage",
				Description:  "Persistent storage for customer documents.",
				Capabilities: map[string]string{},
			},
			wantErr: true,
		},
		{
			name: "empty profile",
			purpose: Purpose{
				Name:        "customer_document_storage",
				Description: "Persistent storage for customer documents.",
				Capabilities: map[string]string{
					"object_storage": "",
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.purpose.Validate()

			if (err != nil) != tt.wantErr {
				t.Fatalf("Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}