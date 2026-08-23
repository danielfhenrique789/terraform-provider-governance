# Terraform Governance Provider

A Terraform provider for applying and monitoring governance definitions maintained in a centralized GitHub repository.

The provider does **not** provision or configure AWS, Azure, GCP, or other infrastructure. It only reads governance definitions and exposes them to Terraform.

## How It Works

The provider connects to a GitHub repository containing a versioned purpose catalog.

```text
GitHub repository
       │
       │ purposes/purposes.yaml
       ▼
Governance Provider
       │
       ▼
governance_purpose
       │
       ▼
Terraform state
```

A purpose defines a set of capabilities and profiles.

For example:

```yaml
version: "1.1"

purposes:
  centralized_audit_logging:
    description: >
      Centralized, protected retention of organization-wide audit and
      configuration logs for security operations, governance, and compliance.

    capabilities:
      object_storage: compliance_archive
      cryptographic-key-management: data_encryption
      audit_logging: centralized
```

The provider resolves the purpose and its inherited capabilities and exposes the result through Terraform.

## Requirements

* Terraform 1.15+
* Go 1.23+ for development
* A GitHub repository containing `purposes/purposes.yaml`
* A GitHub token with permission to read the repository

## Provider Configuration

The provider requires the GitHub repository containing the governance catalog.

```hcl
terraform {
  required_providers {
    governance = {
      source = "danielfhenrique789/governance"
    }
  }
}

provider "governance" {
  repository = "organization/enterprise-capabilities"
}
```

The GitHub token is read from the `GITHUB_TOKEN` environment variable.

```bash
export GITHUB_TOKEN="your-token"
```

The token is not stored in Terraform configuration.

## Governance Purpose

A purpose is represented as a Terraform resource.

```hcl
resource "governance_purpose" "audit_logging" {
  name            = "centralized_audit_logging"
  catalog_version = "1.1"
}
```

The resource exposes the resolved capabilities:

```text
governance_purpose.audit_logging
├── name
├── catalog_version
└── capabilities
```

For example:

```text
capabilities = {
  "audit_logging"                = "centralized"
  "cryptographic-key-management" = "data_encryption"
  "object_storage"               = "compliance_archive"
}
```

## Catalog Versioning

The catalog version is explicitly declared by the Terraform configuration:

```hcl
catalog_version = "1.1"
```

This represents the version of the governance catalog that the team has acknowledged.

When the central catalog is updated, for example from:

```yaml
version: "1.1"
```

to:

```yaml
version: "1.2"
```

the provider compares the configured version with the current catalog version.

If they differ, Terraform produces a warning:

```text
Warning: Purpose catalog has changed

Purpose "centralized_audit_logging" is configured for catalog version "1.1",
but the current catalog is version "1.2".
Review the updated governance requirements.
```

The warning does **not** block the Terraform plan or apply.

This allows teams to continue working while making the governance change visible.

Once the team has reviewed and adopted the new governance requirements, it updates:

```hcl
catalog_version = "1.2"
```

The warning then disappears.

## Purpose Inheritance

Purposes can extend another purpose.

```yaml
version: "1.1"

purposes:
  basic_data_storage:
    description: Basic enterprise data storage.

    capabilities:
      object_storage: standard
      key_management: data_encryption

  compliant_data_storage:
    description: Compliant enterprise data storage.

    extends: basic_data_storage

    capabilities:
      retention: compliance
      immutability: compliance
```

Resolving `compliant_data_storage` produces:

```text
object_storage: standard
key_management: data_encryption
retention: compliance
immutability: compliance
```

The provider rejects:

* Unknown purposes.
* Unknown parent purposes.
* Circular inheritance.
* Capability conflicts between a purpose and its inherited capabilities.

## Repository Structure

The provider expects the governance repository to contain:

```text
purposes/
└── purposes.yaml
```

The catalog contains:

```yaml
version: "1.1"

purposes:
  purpose_name:
    description: Purpose description.

    capabilities:
      capability_name: profile
```

The `version` field is required.

Each purpose must define at least one capability.

## Local Development

Build the provider:

```bash
go build -o terraform-provider-governance
```

Run the tests:

```bash
go test ./...
```

For local Terraform development, configure a development override.

Example `.terraformrc`:

```hcl
provider_installation {
  dev_overrides {
    "registry.terraform.io/danielfhenrique789/governance" = "/path/to/terraform-provider-governance"
  }

  direct {}
}
```

Set the configuration file:

```bash
export TF_CLI_CONFIG_FILE="/path/to/.terraformrc"
```

Build the provider binary in the directory configured by `dev_overrides`.

Terraform will then use the local provider binary instead of downloading a released provider.

## Current MVP Scope

The MVP provides:

* GitHub-backed governance catalog loading.
* YAML catalog parsing and validation.
* Purpose resolution.
* Purpose inheritance.
* Circular inheritance detection.
* Capability conflict detection.
* Terraform `governance_purpose` resource.
* Catalog version awareness.
* Non-blocking warnings when a team's configured catalog version is outdated.
* No direct communication with AWS or other infrastructure providers.

## What the MVP Does Not Do

The provider currently does not:

* Provision infrastructure.
* Configure AWS services.
* Configure IAM policies.
* Apply capability policies automatically.
* Enforce governance changes.
* Automatically update a team's configured catalog version.
* Provide centralized notifications outside Terraform diagnostics.
* Manage the governance catalog itself.

The provider is intentionally limited to **governance definition resolution and drift awareness**.

## Design Principle

Governance changes should be visible without unnecessarily blocking engineering teams.

The provider therefore separates:

```text
Governance detection
        ↓
      Warning
        ↓
Team reviews change
        ↓
Team updates catalog_version
        ↓
Governance acknowledged
```

Enforcement can be added later as a separate, configurable behavior without changing the core catalog model.
