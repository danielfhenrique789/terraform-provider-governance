package purpose

type PurposeCatalog struct {
	Version  string
	Purposes map[string]Purpose
}

type Purpose struct {
	Name         string
	Description  string
	Extends      *string
	Capabilities map[string]string
}

type ResolvedPurpose struct {
	Name         string
	Description  string
	Capabilities map[string]string
}