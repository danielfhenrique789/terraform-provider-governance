package provider

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

var testAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"governance": providerserver.NewProtocol6WithError(New()),
}

func TestAccPurposeResource(t *testing.T) {
	if os.Getenv("TF_ACC") == "" {
		t.Skip("TF_ACC must be set to run acceptance tests")
	}

	resource.UnitTest(t, resource.TestCase{
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccPurposeResourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(
						"governance_purpose.audit_logging",
						"name",
						"centralized_audit_logging",
					),
					resource.TestCheckResourceAttr(
						"governance_purpose.audit_logging",
						"catalog_version",
						"1.0",
					),
				),
			},
		},
	})
}

func testAccPurposeResourceConfig() string {
	return `
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
  name            = "centralized_audit_logging"
  catalog_version = "1.0"
}
`
}
