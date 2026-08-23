package purpose

import "fmt"

func (c PurposeCatalog) Resolve(name string) (ResolvedPurpose, error) {
	if _, exists := c.Purposes[name]; !exists {
		return ResolvedPurpose{}, fmt.Errorf(
			"purpose %q does not exist",
			name,
		)
	}

	capabilities, err := c.resolveCapabilities(name, map[string]bool{})
	if err != nil {
		return ResolvedPurpose{}, err
	}

	p := c.Purposes[name]

	return ResolvedPurpose{
		Name:         p.Name,
		Description:  p.Description,
		Capabilities: capabilities,
	}, nil
}

func (c PurposeCatalog) resolveCapabilities(
	name string,
	resolving map[string]bool,
) (map[string]string, error) {
	p, exists := c.Purposes[name]
	if !exists {
		return nil, fmt.Errorf("purpose %q does not exist", name)
	}

	if resolving[name] {
		return nil, fmt.Errorf(
			"circular purpose inheritance detected at %q",
			name,
		)
	}

	resolving[name] = true
	defer delete(resolving, name)

	result := make(map[string]string)

	if p.Extends != nil {
		parentCapabilities, err := c.resolveCapabilities(
			*p.Extends,
			resolving,
		)
		if err != nil {
			return nil, err
		}

		for capability, profile := range parentCapabilities {
			result[capability] = profile
		}
	}

	for capability, profile := range p.Capabilities {
		if _, exists := result[capability]; exists {
			return nil, fmt.Errorf(
				"purpose %q conflicts with inherited capability %q",
				name,
				capability,
			)
		}

		result[capability] = profile
	}

	return result, nil
}