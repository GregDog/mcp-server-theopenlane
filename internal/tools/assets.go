package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type assetItem struct {
	ID               string   `json:"id"`
	Name             string   `json:"name,omitempty"`
	Identifier       string   `json:"identifier,omitempty"`
	AssetType        string   `json:"asset_type,omitempty"`
	AssetSubtypeName string   `json:"asset_subtype_name,omitempty"`
	CriticalityName  string   `json:"criticality_name,omitempty"`
	EnvironmentName  string   `json:"environment_name,omitempty"`
	Description      string   `json:"description,omitempty"`
	ContainsPii      *bool    `json:"contains_pii,omitempty"`
	Region           string   `json:"region,omitempty"`
	SourceType       string   `json:"source_type,omitempty"`
	Tags             []string `json:"tags,omitempty"`
}

func registerAssets(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_assets_list",
		Title:       "List Openlane assets",
		Description: "List assets in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listAssets)

	addTool(server, &mcp.Tool{
		Name:        "openlane_asset_get",
		Title:       "Get an Openlane asset",
		Description: "Get a single asset by ID.",
		Annotations: readOnly(),
	}, h.getAsset)
}

func (h *handlers) listAssets(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[assetItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetAssets(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[assetItem]{}, openlane.APIError(err)
	}
	items := make([]assetItem, 0, len(resp.Assets.Edges))
	for _, e := range resp.Assets.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListAsset(e.Node))
	}
	return nil, openlane.Page[assetItem]{
		Items:      items,
		NextCursor: resp.Assets.PageInfo.EndCursor,
		HasMore:    resp.Assets.PageInfo.HasNextPage,
		TotalCount: resp.Assets.TotalCount,
	}, nil
}

func (h *handlers) getAsset(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, assetItem, error) {
	if in.ID == "" {
		return nil, assetItem{}, errIDRequired
	}
	resp, err := h.api.GetAssetByID(ctx, in.ID)
	if err != nil {
		return nil, assetItem{}, openlane.APIError(err)
	}
	return nil, mapGetAsset(resp.Asset), nil
}

func mapListAsset(n *graphclient.GetAssets_Assets_Edges_Node) assetItem {
	return assetItem{
		ID:               n.ID,
		Name:             n.Name,
		Identifier:       openlane.Deref(n.Identifier),
		AssetType:        openlane.Format(n.AssetType),
		AssetSubtypeName: openlane.Deref(n.AssetSubtypeName),
		CriticalityName:  openlane.Deref(n.CriticalityName),
		EnvironmentName:  openlane.Deref(n.EnvironmentName),
	}
}

func mapGetAsset(a graphclient.GetAssetByID_Asset) assetItem {
	containsPii := a.ContainsPii
	return assetItem{
		ID:               a.ID,
		Name:             a.Name,
		Identifier:       openlane.Deref(a.Identifier),
		AssetType:        openlane.Format(a.AssetType),
		AssetSubtypeName: openlane.Deref(a.AssetSubtypeName),
		CriticalityName:  openlane.Deref(a.CriticalityName),
		EnvironmentName:  openlane.Deref(a.EnvironmentName),
		Description:      openlane.Deref(a.Description),
		ContainsPii:      containsPii,
		Region:           openlane.Deref(a.Region),
		SourceType:       openlane.Format(a.SourceType),
		Tags:             a.Tags,
	}
}
