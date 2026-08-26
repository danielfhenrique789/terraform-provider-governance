terraform {
  required_providers {
    governance = {
      source = "danielfhenrique789/governance"
    }
  }
}

provider "governance" {
  repository = "danielfhenrique789/enterprise-capabilities"
}

resource "governance_purpose" "audit" {
  name = "centralized_audit_logging"
  catalog_version = "1.0"
}