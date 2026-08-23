package definitions

import (
	"context"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
)

type Repository interface {
	GetPurposeCatalog(ctx context.Context) (purpose.PurposeCatalog, error)
}