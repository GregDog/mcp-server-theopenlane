package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

func TestCreatePlatformRequiresName(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.createPlatform(context.Background(), nil, createPlatformInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreatePlatformSuccess(t *testing.T) {
	api := &fakeAPI{}
	h := &handlers{api: api}
	_, item, err := h.createPlatform(context.Background(), nil, createPlatformInput{
		Name:        "Production API",
		Description: "Customer-facing REST API",
		Region:      "us-east-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if item.Name != "Production API" || api.lastCreatePlatformInput.Name != "Production API" {
		t.Fatalf("unexpected create: %+v input=%+v", item, api.lastCreatePlatformInput)
	}
}

func TestCreatePlatformInvalidStatus(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.createPlatform(context.Background(), nil, createPlatformInput{
		Name:   "Production API",
		Status: "BOGUS",
	})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestCreatePlatformOwnerResolvesUserID(t *testing.T) {
	userID := "USR01J9ABCD111111111111111"
	api := &fakeAPI{
		orgMembers: &graphclient.GetOrgMembersByOrgID{
			OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
				Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
					{Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
						User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
							ID:    userID,
							Email: "owner@example.com",
						},
					}},
				},
			},
		},
	}
	h := &handlers{api: api, organizationID: "org_1"}
	_, _, err := h.createPlatform(context.Background(), nil, createPlatformInput{
		Name:            "Production API",
		TechnicalOwner:  "owner@example.com",
	})
	if err != nil {
		t.Fatal(err)
	}
	if api.lastCreatePlatformInput.TechnicalOwnerUserID == nil || *api.lastCreatePlatformInput.TechnicalOwnerUserID != userID {
		t.Fatalf("expected technical owner user id %q, got %+v", userID, api.lastCreatePlatformInput.TechnicalOwnerUserID)
	}
	if api.lastCreatePlatformInput.TechnicalOwnerGroupID != nil {
		t.Fatalf("expected user id not group id: %+v", api.lastCreatePlatformInput.TechnicalOwnerGroupID)
	}
}

func TestCreatePlatformOwnerResolvesGroupID(t *testing.T) {
	groupID := "GRP01J9ABCD111111111111111"
	api := &fakeAPI{
		groups: &graphclient.GetGroups{
			Groups: graphclient.GetGroups_Groups{
				Edges: []*graphclient.GetGroups_Groups_Edges{
					{Node: &graphclient.GetGroups_Groups_Edges_Node{
						ID:          groupID,
						Name:        "Platform Ops",
						DisplayName: "Platform Ops",
					}},
				},
			},
		},
	}
	h := &handlers{api: api, organizationID: "org_1"}
	_, _, err := h.createPlatform(context.Background(), nil, createPlatformInput{
		Name:            "Production API",
		TechnicalOwner:  "Platform Ops",
	})
	if err != nil {
		t.Fatal(err)
	}
	if api.lastCreatePlatformInput.TechnicalOwnerGroupID == nil || *api.lastCreatePlatformInput.TechnicalOwnerGroupID != groupID {
		t.Fatalf("expected technical owner group id %q, got %+v", groupID, api.lastCreatePlatformInput.TechnicalOwnerGroupID)
	}
}
