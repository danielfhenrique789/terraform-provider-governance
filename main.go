package main

import (
	"context"
	"log"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/provider"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

func main() {
	err := providerserver.Serve(
		context.Background(),
		provider.New,
		providerserver.ServeOpts{
			Address: "registry.terraform.io/danielfhenrique789/governance",
		},
	)

	if err != nil {
		log.Fatal(err)
	}
}