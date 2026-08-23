package definitions

import (
	"context"
	"fmt"
	"strings"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/purpose"
	"github.com/google/go-github/v68/github"
)

type GitHubRepository struct {
	Client  *github.Client
	Owner   string
	Repo    string
	Version string
}

func NewGitHubRepository(
	client *github.Client,
	repository string,
	version string,
) (*GitHubRepository, error) {
	parts := strings.Split(repository, "/")

	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return nil, fmt.Errorf(
			"invalid GitHub repository %q, expected format owner/repository",
			repository,
		)
	}

	return &GitHubRepository{
		Client:  client,
		Owner:   parts[0],
		Repo:    parts[1],
		Version: version,
	}, nil
}

func (r *GitHubRepository) GetPurposeCatalog(
	ctx context.Context,
) (purpose.PurposeCatalog, error) {
	const path = "purposes/purposes.yaml"

	file, _, _, err := r.Client.Repositories.GetContents(
		ctx,
		r.Owner,
		r.Repo,
		path,
		&github.RepositoryContentGetOptions{
			Ref: r.Version,
		},
	)
	if err != nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"getting purpose definitions from GitHub: %w",
			err,
		)
	}

	if file == nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"purpose definitions file is not a file",
		)
	}

	content, err := file.GetContent()
	if err != nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"reading purpose definitions: %w",
			err,
		)
	}

	catalog, err := LoadPurposeCatalog([]byte(content))
	if err != nil {
		return purpose.PurposeCatalog{}, fmt.Errorf(
			"parsing purpose definitions: %w",
			err,
		)
	}

	return catalog, nil
}

var _ Repository = (*GitHubRepository)(nil)
