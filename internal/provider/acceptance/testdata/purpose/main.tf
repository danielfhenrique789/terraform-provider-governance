terraform {
  required_providers {
    governance = {
      source = "registry.terraform.io/danielfhenrique789/governance"
    }
  }
}

provider "governance" {
  repository = "danielfhenrique789/enterprise-capabilities"
}

resource "governance_purpose" "audit_logging" {
  name = "centralized_audit_logging"
  catalog_version = "1.0"
}