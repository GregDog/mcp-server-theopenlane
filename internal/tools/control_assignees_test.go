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
	managed := true
	api := &fakeAPI{
		orgMembers: &graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								ID:          userID,
								Email:       email,
								DisplayName: "Alice Example",
							},
						},
					},
				},
			},
		},
		groupsByID: map[string]graphclient.GetGroupByID_Group{
			ownerGroup: {
				ID:          ownerGroup,
				Name:        "alice - " + userID,
				DisplayName: "alice",
				IsManaged:   &managed,
			},
			delegateGroup: {
				ID:          delegateGroup,
				Name:        "delegate",
				DisplayName: "delegate",
			},
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

func TestUserIDFromManagedGroupName(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU04"
	got := userIDFromManagedGroupName("Greg - " + userID)
	if got != userID {
		t.Fatalf("got %q want %q", got, userID)
	}
	if userIDFromManagedGroupName("Security Team") != "" {
		t.Fatal("expected empty for non-managed group name")
	}
}

func TestResolveGroupAssigneeSummaryIncludesUser(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU05"
	groupID := "01GROUPGROUPGROUPGROUPGR07"
	email := "greg@example.com"
	api := &fakeAPI{
		group: &graphclient.GetGroupByID{
			Group: graphclient.GetGroupByID_Group{
				ID:          groupID,
				Name:        "Greg - " + userID,
				DisplayName: "Greg",
				IsManaged:   boolPtr(true),
			},
		},
		orgMembers: &graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								ID:          userID,
								DisplayName: "Greg Knell",
								Email:       email,
							},
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api, organizationID: "org_1"}
	got, err := h.resolveGroupAssigneeSummary(context.Background(), groupID, nil)
	if err != nil {
		t.Fatal(err)
	}
	if got.UserID != userID || got.UserEmail != email || got.UserDisplayName != "Greg Knell" {
		t.Fatalf("unexpected summary: %+v", got)
	}
}

func TestGetControlEnrichesDelegate(t *testing.T) {
	userID := "01USERUSERUSERUSERUSERU06"
	groupID := "01GROUPGROUPGROUPGROUPGR08"
	title := "Security awareness"
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:             "ctrl_1",
				RefCode:        "12.6.2",
				Title:          &title,
				DelegateID:     &groupID,
				ControlOwnerID: &groupID,
			},
		},
		group: &graphclient.GetGroupByID{
			Group: graphclient.GetGroupByID_Group{
				ID:          groupID,
				Name:        "alice - " + userID,
				DisplayName: "alice",
				IsManaged:   boolPtr(true),
			},
		},
		orgMembers: &graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{
						Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								ID:          userID,
								DisplayName: "Alice Example",
								Email:       "alice@example.com",
							},
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api, organizationID: "org_1"}
	_, item, err := h.getControl(context.Background(), nil, getInput{ID: "ctrl_1"})
	if err != nil {
		t.Fatal(err)
	}
	if item.Delegate == nil || item.Delegate.UserEmail != "alice@example.com" {
		t.Fatalf("delegate: %+v", item.Delegate)
	}
	if item.ControlOwner == nil || item.ControlOwner.UserID != userID {
		t.Fatalf("control owner: %+v", item.ControlOwner)
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
