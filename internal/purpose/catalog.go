package purpose

import "fmt"

func (c PurposeCatalog) Validate() error {
	if c.Version == "" {
		return fmt.Errorf("catalog version is required")
	}

	if len(c.Purposes) == 0 {
		return fmt.Errorf("catalog must contain at least one purpose")
	}

	for name, p := range c.Purposes {
		if name == "" {
			return fmt.Errorf("catalog contains a purpose with an empty name")
		}

		if p.Name != name {
			return fmt.Errorf(
				"purpose name %q does not match catalog key %q",
				p.Name,
				name,
			)
		}

		if err := p.Validate(); err != nil {
			return fmt.Errorf("purpose %q is invalid: %w", name, err)
		}
	}

	if err := c.ValidateInheritance(); err != nil {
		return fmt.Errorf("invalid purpose inheritance: %w", err)
	}

	return nil
}