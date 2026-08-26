package datasources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
)

func TestCapabilityDataSourceMetadata(t *testing.T) {
	dataSource := CapabilityDataSource{}

	var response datasource.MetadataResponse

	dataSource.Metadata(
		context.Background(),
		datasource.MetadataRequest{
			ProviderTypeName: "governance",
		},
		&response,
	)

	if response.TypeName != "governance_capability" {
		t.Fatalf(
			"unexpected data source type name: %q",
			response.TypeName,
		)
	}
}
