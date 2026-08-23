package purpose

import "fmt"

func (p Purpose) Validate() error {
	if p.Name == "" {
		return fmt.Errorf("purpose name is required")
	}

	if p.Description == "" {
		return fmt.Errorf("purpose description is required")
	}

	if len(p.Capabilities) == 0 {
		return fmt.Errorf("purpose must define at least one capability")
	}

	for capability, profile := range p.Capabilities {
		if capability == "" {
			return fmt.Errorf("purpose contains an empty capability name")
		}

		if profile == "" {
			return fmt.Errorf("purpose capability %q has an empty profile", capability)
		}
	}

	return nil
}