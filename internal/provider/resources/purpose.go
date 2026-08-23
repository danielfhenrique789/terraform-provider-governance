package resources

import (
	"context"
	"fmt"

	"github.com/danielfhenrique789/terraform-provider-governance/internal/definitions"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type PurposeResource struct {
	repository definitions.Repository
}

type PurposeResourceModel struct {
	Name           types.String `tfsdk:"name"`
	Capabilities   types.Map    `tfsdk:"capabilities"`
	CatalogVersion types.String `tfsdk:"catalog_version"`
}

func (r PurposeResource) Metadata(
	ctx context.Context,
	req resource.MetadataRequest,
	resp *resource.MetadataResponse,
) {
	resp.TypeName = req.ProviderTypeName + "_purpose"
}

func (r PurposeResource) Schema(
	ctx context.Context,
	req resource.SchemaRequest,
	resp *resource.SchemaResponse,
) {
	resp.Schema = schema.Schema{
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"capabilities": schema.MapAttribute{
				Computed:    true,
				ElementType: types.StringType,
			},
			"catalog_version": schema.StringAttribute{
				Required: true,
			},
		},
	}
}

func (r *PurposeResource) Configure(
	ctx context.Context,
	req resource.ConfigureRequest,
	resp *resource.ConfigureResponse,
) {
	if req.ProviderData == nil {
		return
	}

	repository, ok := req.ProviderData.(definitions.Repository)

	if !ok {
		resp.Diagnostics.AddError(
			"Invalid provider configuration",
			"Could not retrieve the definition repository.",
		)
		return
	}

	r.repository = repository
}

func capabilityValues(
	capabilities map[string]string,
) map[string]attr.Value {
	values := make(map[string]attr.Value, len(capabilities))

	for capability, profile := range capabilities {
		values[capability] = types.StringValue(profile)
	}

	return values
}

func (r PurposeResource) Create(
	ctx context.Context,
	req resource.CreateRequest,
	resp *resource.CreateResponse,
) {
	if r.repository == nil {
		resp.Diagnostics.AddError(
			"Missing definition repository",
			"The definition repository was not configured.",
		)
		return
	}

	var data PurposeResourceModel

	resp.Diagnostics.Append(req.Plan.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := r.repository.GetPurposeCatalog(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to load purpose definitions",
			err.Error(),
		)
		return
	}

	resolved, err := catalog.Resolve(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Purpose validation failed",
			err.Error(),
		)
		return
	}

	data.Capabilities = types.MapValueMust(
		types.StringType,
		capabilityValues(resolved.Capabilities),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r PurposeResource) Read(
	ctx context.Context,
	req resource.ReadRequest,
	resp *resource.ReadResponse,
) {
	if r.repository == nil {
		resp.Diagnostics.AddError(
			"Missing definition repository",
			"The definition repository was not configured.",
		)
		return
	}

	var data PurposeResourceModel

	resp.Diagnostics.Append(req.State.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	catalog, err := r.repository.GetPurposeCatalog(ctx)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to load purpose definitions",
			err.Error(),
		)
		return
	}

	resolved, err := catalog.Resolve(data.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Purpose validation failed",
			err.Error(),
		)
		return
	}

	data.Capabilities = types.MapValueMust(
		types.StringType,
		capabilityValues(resolved.Capabilities),
	)

	resp.Diagnostics.Append(resp.State.Set(ctx, &data)...)
}

func (r PurposeResource) ModifyPlan(
	ctx context.Context,
	req resource.ModifyPlanRequest,
	resp *resource.ModifyPlanResponse,
) {
	if r.repository == nil {
		return
	}

	var data PurposeResourceModel

	resp.Diagnostics.Append(
		req.Plan.Get(ctx, &data)...,
	)

	if resp.Diagnostics.HasError() {
		return
	}

	if data.CatalogVersion.IsNull() ||
		data.CatalogVersion.IsUnknown() {
		return
	}

	catalog, err := r.repository.GetPurposeCatalog(ctx)
	if err != nil {
		resp.Diagnostics.AddWarning(
			"Unable to check purpose catalog version",
			err.Error(),
		)
		return
	}

	if data.CatalogVersion.ValueString() != catalog.Version {
		resp.Diagnostics.AddWarning(
			"Purpose catalog has changed",
			fmt.Sprintf(
				"Purpose %q is configured for catalog version %q, but the current catalog is version %q. Review the updated governance requirements.",
				data.Name.ValueString(),
				data.CatalogVersion.ValueString(),
				catalog.Version,
			),
		)
	}
}

func (r PurposeResource) Update(
	ctx context.Context,
	req resource.UpdateRequest,
	resp *resource.UpdateResponse,
) {
}

func (r PurposeResource) Delete(
	ctx context.Context,
	req resource.DeleteRequest,
	resp *resource.DeleteResponse,
) {
}

func NewPurposeResource() resource.Resource {
	return &PurposeResource{}
}
