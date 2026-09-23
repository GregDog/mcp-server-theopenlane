package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

func TestResolveControlAssigneeGroupIDUsesGroupID(t *testing.T) {
	groupID := "01GROUPGROUPGROUPGROUPGR01"
	api := &fakeAPI{
		group: &graphclient.GetGroupByID{
			Group: graphclient.GetGroupByID_Group{ID: groupID, Name: "Security"},
		},
	}
	h := &handlers{api: api}
	got, err := h.resolveControlAssigneeGroupID(context.Background(), groupID)
	if err != nil {
		t.Fatal(err)
	}
	if got != groupID {
		t.Fatalf("group id: got %q want %q", got, groupID)
	}
}

func TestResolveControlAssigneeGroupIDResolvesUserToManagedGroup(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU01"
	groupID := "01GROUPGROUPGROUPGROUPGR02"
	api := &fakeAPI{
		group: &graphclient.GetGroupByID{},
		groups: &graphclient.GetGroups{
			Groups: graphclient.GetGroups_Groups{
				Edges: []*graphclient.GetGroups_Groups_Edges{
					{
						Node: &graphclient.GetGroups_Groups_Edges_Node{
							ID:   groupID,
							Name: "alice - " + userID,
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api}
	got, err := h.resolveControlAssigneeGroupID(context.Background(), userID)
	if err != nil {
		t.Fatal(err)
	}
	if got != groupID {
		t.Fatalf("group id: got %q want %q", got, groupID)
	}
}

func TestUpdateControlResolvesOwnerAndDelegate(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU02"
	ownerGroup := "01GROUPGROUPGROUPGROUPGR03"
	delegateGroup := "01GROUPGROUPGROUPGROUPGR04"
	email := "alice@example.com"
	api := &fakeAPI{
		orgMembers: &graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								ID:    userID,
								Email: email,
							},
						},
					},
				},
			},
		},
		group: &graphclient.GetGroupByID{
			Group: graphclient.GetGroupByID_Group{ID: delegateGroup, Name: "delegate"},
		},
		groups: &graphclient.GetGroups{
			Groups: graphclient.GetGroups_Groups{
				Edges: []*graphclient.GetGroups_Groups_Edges{
					{
						Node: &graphclient.GetGroups_Groups_Edges_Node{
							ID:   ownerGroup,
							Name: "alice - " + userID,
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api, allowWrite: true, organizationID: "org_1"}

	_, _, err := h.updateControl(context.Background(), nil, updateControlInput{
		ID:         "ctrl_1",
		OwnerID:    email,
		DelegateID: delegateGroup,
	})
	if err != nil {
		t.Fatal(err)
	}
	if api.lastUpdateControlInput.ControlOwnerID == nil || *api.lastUpdateControlInput.ControlOwnerID != ownerGroup {
		t.Fatalf("control owner: %+v", api.lastUpdateControlInput.ControlOwnerID)
	}
	if api.lastUpdateControlInput.DelegateID == nil || *api.lastUpdateControlInput.DelegateID != delegateGroup {
		t.Fatalf("delegate: %+v", api.lastUpdateControlInput.DelegateID)
	}
}

func TestPickPersonalManagedGroupPrefersUserSuffix(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU03"
	want := "01GROUPGROUPGROUPGROUPGR05"
	got, err := pickPersonalManagedGroup(userID, []graphclient.GetGroups_Groups_Edges_Node{
		{ID: "01GROUPGROUPGROUPGROUPGR06", Name: "Admin User - other"},
		{ID: want, Name: "alice - " + userID},
	})
	if err != nil {
		t.Fatal(err)
	}
	if got != want {
		t.Fatalf("got %q want %q", got, want)
	}
}
