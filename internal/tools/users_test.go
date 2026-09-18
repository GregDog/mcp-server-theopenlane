package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

func TestListUsersReturnsOrgMembers(t *testing.T) {
	email := "greg.knell@nomupay.com"
	h := &handlers{
		organizationID: "01ORG",
		api: &fakeAPI{
			orgMembers: &graphclient.GetOrgMembersByOrgID{
				OrgMemberships: graphclient.GetOrgMembersByOrgID_OrgMemberships{
					Edges: []*graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges{
						{Node: &graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node{
							User: graphclient.GetOrgMembersByOrgID_OrgMemberships_Edges_Node_User{
								ID:          "01USER",
								Email:       email,
								DisplayName: "Greg",
							},
						}},
					},
				},
			},
		},
	}
	_, page, err := h.listUsers(context.Background(), nil, userListInput{Email: email})
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if len(page.Items) != 1 || page.Items[0].Email != email {
		t.Fatalf("unexpected page: %+v", page)
	}
}
