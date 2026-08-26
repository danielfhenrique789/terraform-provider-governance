package provider

import (
	"context"
	"os"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/definitions"
	"github.com/danielfhenrique789/terraform-provider-governance/internal/provider/datasources"
	"github.com/danielfhenrique789/terraform-provider-governance/internal/provider/resources"
	"github.com/google/go-github/v68/github"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/provider"
	"github.com/hashicorp/terraform-plugin-framework/provider/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type GovernanceProvider struct{}

type GovernanceProviderModel struct {
	Repository types.String `tfsdk:"repository"`
}

func (p *GovernanceProvider) Metadata(
	ctx context.Context,
	req provider.MetadataRequest,
	resp *provider.MetadataResponse,
) {
	resp.TypeName = "governance"
}

func (p *GovernanceProvider) Configure(
	ctx context.Context,
	req provider.ConfigureRequest,
	resp *provider.ConfigureResponse,
) {
	var data GovernanceProviderModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.Repository.IsNull() || data.Repository.ValueString() == "" {
		resp.Diagnostics.AddError(
			"Missing repository",
			"The repository attribute must be configured.",
		)
		return
	}

	token := os.Getenv("GITHUB_TOKEN")

	if token == "" {
		resp.Diagnostics.AddError(
			"Missing GitHub token",
			"The GITHUB_TOKEN environment variable must be configured.",
		)
		return
	}

	client := github.NewClient(nil).WithAuthToken(token)

	repository, err := definitions.NewGitHubRepository(
		client,
		data.Repository.ValueString(),
		"main",
	)
	if err != nil {
		resp.Diagnostics.AddError(
			"Invalid GitHub repository",
			err.Error(),
		)
		return
	}

	resp.ResourceData = repository
	resp.DataSourceData = repository
}

func (p *GovernanceProvider) Resources(
	ctx context.Context,
) []func() resource.Resource {
	return []func() resource.Resource{
		resources.NewPurposeResource,
	}
}

func (p *GovernanceProvider) DataSources(
	ctx context.Context,
) []func() datasource.DataSource {
	return []func() datasource.DataSource{
		datasources.NewPurposeDataSource,
		datasources.NewCapabilityDataSource,
		datasources.NewCapabilityProfileDataSource,
	}
}

func (p *GovernanceProvider) Schema(
	ctx context.Context,
	req provider.SchemaRequest,
	resp *provider.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"repository": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func New() provider.Provider {
	return &GovernanceProvider{}
}
