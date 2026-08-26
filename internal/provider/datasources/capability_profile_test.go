package datasources

import (
	"strings"
	"testing"
)

func TestResolveCapabilityProfile(t *testing.T) {
	tests := map[string]struct {
		capability        string
		profile           string
		expectedVersion   string
		expectedService   string
		expectedSubstring string
	}{
		"object storage compliance archive": {
			capability:        "object-storage",
			profile:           "compliance_archive",
			expectedVersion:   "1.0",
			expectedService:   "s3",
			expectedSubstring: `"mode":"COMPLIANCE"`,
		},
		"cryptographic key management data encryption": {
			capability:        "cryptographic-key-management",
			profile:           "data_encryption",
			expectedVersion:   "1.0",
			expectedService:   "kms",
			expectedSubstring: `"usage":"encrypt_decrypt"`,
		},
		"audit logging centralized": {
			capability:        "audit-logging",
			profile:           "centralized",
			expectedVersion:   "1.0",
			expectedService:   "audit_logging",
			expectedSubstring: `"scope":"organization"`,
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			version, service, configuration, err :=
				resolveCapabilityProfile(tt.capability, tt.profile)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if version != tt.expectedVersion {
				t.Errorf(
					"version = %q, want %q",
					version,
					tt.expectedVersion,
				)
			}

			if service != tt.expectedService {
				t.Errorf(
					"service = %q, want %q",
					service,
					tt.expectedService,
				)
			}

			if !strings.Contains(configuration, tt.expectedSubstring) {
				t.Errorf(
					"configuration does not contain %q: %s",
					tt.expectedSubstring,
					configuration,
				)
			}
		})
	}
}

func TestResolveCapabilityProfileUnknownCapability(t *testing.T) {
	_, _, _, err := resolveCapabilityProfile(
		"unknown-capability",
		"standard",
	)

	if err == nil {
		t.Fatal("expected error for unknown capability")
	}
}

func TestResolveCapabilityProfileUnknownProfile(t *testing.T) {
	_, _, _, err := resolveCapabilityProfile(
		"object-storage",
		"unknown-profile",
	)

	if err == nil {
		t.Fatal("expected error for unknown profile")
	}
}
