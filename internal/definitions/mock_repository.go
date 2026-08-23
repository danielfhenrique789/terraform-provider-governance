package definitions

import (
	"context"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
)

type MockRepository struct {
	Catalog purpose.PurposeCatalog
}

func (r *MockRepository) GetPurposeCatalog(
	ctx context.Context,
) (purpose.PurposeCatalog, error) {
	return r.Catalog, nil
}

var _ Repository = (*MockRepository)(nil)