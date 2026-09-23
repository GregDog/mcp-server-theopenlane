package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"
)

func TestCreateMappedControlRequiresEndpoints(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, _, err := h.createMappedControl(context.Background(), nil, createMappedControlInput{
		MappingType: "EQUAL",
		ToControlIDs: []string{"ctrl_2"},
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreateMappedControlMapsInput(t *testing.T) {
	api := &fakeAPI{}
	h := &handlers{api: api, allowWrite: true}
	_, item, err := h.createMappedControl(context.Background(), nil, createMappedControlInput{
		MappingType:    "SUBSET",
		FromControlIDs: []string{"ctrl_1"},
		ToControlIDs:   []string{"ctrl_2"},
		Source:         "MANUAL",
	})
	if err != nil {
		t.Fatal(err)
	}
	if api.lastCreateMappedControlInput.MappingType == nil || *api.lastCreateMappedControlInput.MappingType != enums.MappingTypeSubset {
		t.Fatalf("mapping type: %+v", api.lastCreateMappedControlInput.MappingType)
	}
	if len(api.lastCreateMappedControlInput.FromControlIDs) != 1 || api.lastCreateMappedControlInput.FromControlIDs[0] != "ctrl_1" {
		t.Fatalf("from controls: %+v", api.lastCreateMappedControlInput.FromControlIDs)
	}
	if item.ID != "map_1" || item.MappingType != "SUBSET" {
		t.Fatalf("item: %+v", item)
	}
}

func TestListMappedControlsFiltersByControlID(t *testing.T) {
	controlID := "ctrl_1"
	api := &fakeAPI{
		mappedControls: &graphclient.GetMappedControls{
			MappedControls: graphclient.GetMappedControls_MappedControls{
				TotalCount: 1,
				PageInfo: graphclient.GetMappedControls_MappedControls_PageInfo{
					HasNextPage: false,
				},
				Edges: []*graphclient.GetMappedControls_MappedControls_Edges{
					{
						Node: &graphclient.GetMappedControls_MappedControls_Edges_Node{
							ID:          "map_1",
							MappingType: enums.MappingTypeEqual,
							FromControls: graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls{
								Edges: []*graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls_Edges{
									{Node: &graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls_Edges_Node{ID: controlID, RefCode: "12.6.1"}},
								},
							},
							ToControls: graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls{
								Edges: []*graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls_Edges{
									{Node: &graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls_Edges_Node{ID: "ctrl_2", RefCode: "12.6.2"}},
								},
							},
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api}
	_, page, err := h.listMappedControls(context.Background(), nil, mappedControlListInput{
		Limit:     10,
		ControlID: controlID,
	})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 1 || len(page.Items) != 1 {
		t.Fatalf("page: %+v", page)
	}
	if len(page.Items[0].FromControls) != 1 || page.Items[0].FromControls[0].RefCode != "12.6.1" {
		t.Fatalf("from controls: %+v", page.Items[0].FromControls)
	}
}

func TestGetControlIncludesRelatedControls(t *testing.T) {
	controlID := "ctrl_from"
	otherID := "ctrl_to"
	api := &fakeAPI{
		controls: emptyControls(),
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      controlID,
				RefCode: "12.6.1",
			},
		},
		mappedControls: &graphclient.GetMappedControls{
			MappedControls: graphclient.GetMappedControls_MappedControls{
				TotalCount: 1,
				Edges: []*graphclient.GetMappedControls_MappedControls_Edges{
					{
						Node: &graphclient.GetMappedControls_MappedControls_Edges_Node{
							ID:          "map_1",
							MappingType: enums.MappingTypeIntersect,
							FromControls: graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls{
								Edges: []*graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls_Edges{
									{Node: &graphclient.GetMappedControls_MappedControls_Edges_Node_FromControls_Edges_Node{ID: controlID, RefCode: "12.6.1"}},
								},
							},
							ToControls: graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls{
								Edges: []*graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls_Edges{
									{Node: &graphclient.GetMappedControls_MappedControls_Edges_Node_ToControls_Edges_Node{ID: otherID, RefCode: "12.6.2"}},
								},
							},
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api}
	_, item, err := h.getControl(context.Background(), nil, getInput{ID: controlID})
	if err != nil {
		t.Fatal(err)
	}
	if item.RelatedControls == nil || item.RelatedControls.Count != 1 || len(item.RelatedControls.Items) != 1 {
		t.Fatalf("related controls: %+v", item.RelatedControls)
	}
	got := item.RelatedControls.Items[0]
	if got.Direction != "outgoing" || got.ID != otherID || got.RefCode != "12.6.2" {
		t.Fatalf("related control: %+v", got)
	}
}

func TestDeleteMappedControl(t *testing.T) {
	api := &fakeAPI{}
	h := &handlers{api: api, allowWrite: true}
	_, result, err := h.deleteMappedControl(context.Background(), nil, getInput{ID: "map_1"})
	if err != nil {
		t.Fatal(err)
	}
	if result.DeletedID != "map_1" || api.deletedID != "map_1" {
		t.Fatalf("delete result: %+v deleted=%q", result, api.deletedID)
	}
}
