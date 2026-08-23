package purpose

import "fmt"

func (c PurposeCatalog) ValidateInheritance() error {
	for name := range c.Purposes {
		if err := c.validateInheritanceChain(name, map[string]bool{}); err != nil {
			return err
		}
	}

	return nil
}

func (c PurposeCatalog) validateInheritanceChain(
	name string,
	visited map[string]bool,
) error {
	p, exists := c.Purposes[name]
	if !exists {
		return fmt.Errorf("purpose %q does not exist", name)
	}

	if visited[name] {
		return fmt.Errorf("circular purpose inheritance detected at %q", name)
	}

	if p.Extends == nil {
		return nil
	}

	parent := *p.Extends

	visited[name] = true
	defer delete(visited, name)

	if _, exists := c.Purposes[parent]; !exists {
		return fmt.Errorf(
			"purpose %q extends unknown purpose %q",
			name,
			parent,
		)
	}

	return c.validateInheritanceChain(parent, visited)
}