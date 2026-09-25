package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func TestListPlatforms(t *testing.T) {
	api := &fakeAPI{
		platforms: &graphclient.GetPlatforms{
			Platforms: graphclient.GetPlatforms_Platforms{
				TotalCount: 1,
				PageInfo:   graphclient.GetPlatforms_Platforms_PageInfo{},
				Edges: []*graphclient.GetPlatforms_Platforms_Edges{
					{Node: &graphclient.GetPlatforms_Platforms_Edges_Node{
						ID:        "plt_1",
						DisplayID: "PLT-001",
						Name:      "Production API",
						Status:    enums.PlatformStatusActive,
					}},
				},
			},
		},
	}
	h := &handlers{api: api}
	_, page, err := h.listPlatforms(context.Background(), nil, platformListInput{Limit: 10, Status: "ACTIVE"})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 1 || len(page.Items) != 1 || page.Items[0].Name != "Production API" {
		t.Fatalf("unexpected page: %+v", page)
	}
	if api.lastPlatformWhere == nil || api.lastPlatformWhere.Status == nil {
		t.Fatalf("expected status filter: %+v", api.lastPlatformWhere)
	}
}

func TestGetPlatformRequiresID(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.getPlatform(context.Background(), nil, getInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetPlatform(t *testing.T) {
	api := &fakeAPI{
		platformDetail: &openlane.PlatformDetail{
			ID:            "plt_1",
			DisplayID:     "PLT-001",
			Name:          "Production API",
			Status:        enums.PlatformStatusActive,
			BusinessOwner: openlane.PlatformOwnerRole{Kind: "name", Name: "Platform Team"},
		},
	}
	h := &handlers{api: api}
	_, item, err := h.getPlatform(context.Background(), nil, getInput{ID: "plt_1"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "Production API" || item.BusinessOwner == nil || item.BusinessOwner.Name != "Platform Team" {
		t.Fatalf("unexpected item: %+v", item)
	}
}
