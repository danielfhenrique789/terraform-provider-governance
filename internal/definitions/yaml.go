package definitions

import (
	"fmt"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
	"gopkg.in/yaml.v3"
)

type purposeCatalogYAML struct {
	Version  string                  `yaml:"version"`
	Purposes map[string]purposeYAML `yaml:"purposes"`
}

type purposeYAML struct {
	Description  string            `yaml:"description"`
	Extends      *string           `yaml:"extends"`
	Capabilities map[string]string `yaml:"capabilities"`
}

func LoadPurposeCatalog(data []byte) (purpose.PurposeCatalog, error) {
	var raw purposeCatalogYAML

	if err := yaml.Unmarshal(data, &raw); err != nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"decode purpose catalog: %w",
			err,
		)
	}

	catalog := purpose.PurposeCatalog{
		Version:  raw.Version,
		Purposes: make(map[string]purpose.Purpose, len(raw.Purposes)),
	}

	for name, rawPurpose := range raw.Purposes {
		catalog.Purposes[name] = purpose.Purpose{
			Name:         name,
			Description:  rawPurpose.Description,
			Extends:      rawPurpose.Extends,
			Capabilities: rawPurpose.Capabilities,
		}
	}

	if err := catalog.Validate(); err != nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"validate purpose catalog: %w",
			err,
		)
	}

	return catalog, nil
}