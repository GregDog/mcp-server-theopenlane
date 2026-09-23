//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	gqlclient "github.com/theopenlane/go-client"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}
	api, err := openlane.New(cfg)
	if err != nil {
		fmt.Println("client:", err)
		os.Exit(1)
	}
	opts := []gqlclient.ClientOption{
		gqlclient.WithAPIToken(cfg.APIToken),
		gqlclient.WithBaseURL(cfg.BaseURL),
	}
	if cfg.OrganizationID != "" {
		opts = append(opts, gqlclient.WithInterceptors(gqlclient.WithOrganizationHeader(cfg.OrganizationID)))
	}
	raw, err := gqlclient.New(opts...)
	if err != nil {
		fmt.Println("raw client:", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	fmt.Printf("org_id=%s base=%s\n", cfg.OrganizationID, cfg.BaseURL)

	first := int64(50)
	users, err := api.GetUsers(ctx, &first, nil, nil)
	if err != nil {
		fmt.Println("GetUsers error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetUsers total=%d edges=%d\n", users.Users.TotalCount, len(users.Users.Edges))
		for _, e := range users.Users.Edges {
			if e != nil && e.Node != nil {
				fmt.Printf("  user id=%s email=%s name=%s\n", e.Node.ID, e.Node.Email, e.Node.DisplayName)
			}
		}
	}

	groups, err := api.GetGroups(ctx, &first, nil, nil)
	if err != nil {
		fmt.Println("GetGroups error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetGroups total=%d edges=%d\n", groups.Groups.TotalCount, len(groups.Groups.Edges))
	}

	self, err := raw.GetSelf(ctx)
	if err != nil {
		fmt.Println("GetSelf error:", openlane.Redact(err.Error()))
	} else {
		defaultOrg := ""
		if self.Self.Setting.DefaultOrg != nil {
			defaultOrg = self.Self.Setting.DefaultOrg.ID
		}
		fmt.Printf("GetSelf id=%s email=%s name=%s default_org=%s\n",
			self.Self.ID, self.Self.Email, self.Self.DisplayName, defaultOrg)
	}

	orgID := cfg.OrganizationID
	members, err := raw.GetOrgMembersByOrgID(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
	})
	if err != nil {
		fmt.Println("GetOrgMembersByOrgID error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetOrgMembersByOrgID edges=%d\n", len(members.OrgMemberships.Edges))
		for _, e := range members.OrgMemberships.Edges {
			if e == nil || e.Node == nil {
				continue
			}
			n := e.Node
			if n.User.ID != "" {
				fmt.Printf("  member user_id=%s email=%s name=%s role=%v\n", n.User.ID, n.User.Email, n.User.DisplayName, n.Role)
			} else {
				fmt.Printf("  member user_id=%s role=%v (no user nested)\n", n.UserID, n.Role)
			}
		}
	}

	allMembers, err := raw.GetOrgMemberships(ctx, &first, nil, nil, nil, nil, nil)
	if err != nil {
		fmt.Println("GetOrgMemberships error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetOrgMemberships total=%d edges=%d\n", allMembers.OrgMemberships.TotalCount, len(allMembers.OrgMemberships.Edges))
	}

	email := os.Getenv("OPENLANE_PROBE_USER_EMAIL")
	if email == "" {
		email = "user@example.com"
	}
	members, err = raw.GetOrgMembersByOrgID(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
		HasUserWith:    []*graphclient.UserWhereInput{{EmailContainsFold: &email}},
	})
	if err != nil {
		fmt.Println("GetOrgMembersByOrgID email filter error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetOrgMembersByOrgID email filter edges=%d\n", len(members.OrgMemberships.Edges))
	}

	uid := "01M1EN3HRS4XJKQNX4R3DJHK0H"
	byID, err := api.GetUserByID(ctx, uid)
	if err != nil {
		fmt.Println("GetUserByID error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetUserByID id=%s email=%s\n", byID.User.ID, byID.User.Email)
	}

	byMember, err := raw.GetOrgMembersByOrgID(ctx, &graphclient.OrgMembershipWhereInput{
		OrganizationID: &orgID,
		UserID:         &uid,
	})
	if err != nil {
		fmt.Println("GetOrgMembersByOrgID userID filter error:", openlane.Redact(err.Error()))
	} else {
		fmt.Printf("GetOrgMembersByOrgID userID filter edges=%d\n", len(byMember.OrgMemberships.Edges))
		if len(byMember.OrgMemberships.Edges) > 0 && byMember.OrgMemberships.Edges[0].Node != nil {
			u := byMember.OrgMemberships.Edges[0].Node.User
			fmt.Printf("  via membership email=%s name=%s\n", u.Email, u.DisplayName)
		}
	}
}
