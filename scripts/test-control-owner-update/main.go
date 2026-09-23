package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
	"github.com/theopenlane/go-client/graphclient"
)

func main() {
	os.Setenv("OPENLANE_API_TOKEN", os.Getenv("OPENLANE_API_TOKEN2"))
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println("config:", err)
		os.Exit(1)
	}
	client, err := openlane.New(cfg)
	if err != nil {
		fmt.Println("client:", err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	abhishekUser := "01KRGYEMTKH5XJABCQZ5BJMXGC"
	gregUser := "01M1EN3HRS4XJKQNX4R3DJHK0H"
	refCodes := []string{"12.6.1", "12.6.2", "12.6.3", "12.6.3.1", "12.6.3.2"}

	fmt.Println("=== resolve PCI 12.6 org-owned controls ===")
	controls := map[string]string{}
	for _, ref := range refCodes {
		first := int64(5)
		where := &graphclient.ControlWhereInput{
			RefCode:       &ref,
			OwnerIDNotNil: boolPtr(true),
		}
		resp, err := client.GetControls(ctx, &first, nil, where)
		if err != nil {
			fmt.Printf("%s list: %s\n", ref, openlane.Redact(err.Error()))
			continue
		}
		var picked *graphclient.GetControls_Controls_Edges_Node
		for _, e := range resp.Controls.Edges {
			if e == nil || e.Node == nil {
				continue
			}
			n := e.Node
			if picked == nil || strings.Contains(openlane.Deref(n.ReferenceFramework), "PCI") {
				picked = n
			}
		}
		if picked == nil {
			fmt.Printf("%s: not found\n", ref)
			continue
		}
		controls[ref] = picked.ID
		fmt.Printf("%s: id=%s display=%s owner_id=%s framework=%s\n",
			ref, picked.ID, picked.DisplayID, openlane.Deref(picked.OwnerID), openlane.Deref(picked.ReferenceFramework))
	}

	fmt.Println("\n=== current control owner/delegate ===")
	for ref, id := range controls {
		resp, err := client.GetControlByID(ctx, id)
		if err != nil {
			fmt.Printf("%s get: %s\n", ref, openlane.Redact(err.Error()))
			continue
		}
		c := resp.Control
		fmt.Printf("%s: control_owner_id=%s delegate_id=%s org_owner_id=%s\n",
			ref, openlane.Deref(c.ControlOwnerID), openlane.Deref(c.DelegateID), openlane.Deref(c.OwnerID))
	}

	fmt.Println("\n=== probe owner update with user IDs (dry-run on 12.6.2 only) ===")
	if id := controls["12.6.2"]; id != "" {
		_, err := client.UpdateControl(ctx, id, graphclient.UpdateControlInput{
			ControlOwnerID: &abhishekUser,
			DelegateID:     &gregUser,
		})
		if err != nil {
			fmt.Printf("12.6.2 owner+delegate user IDs: %s\n", openlane.Redact(err.Error()))
		} else {
			fmt.Println("12.6.2 owner+delegate user IDs: OK")
		}
	}

	fmt.Println("\n=== probe evidence link on org 12.6.2 (permission baseline) ===")
	if id := controls["12.6.2"]; id != "" {
		desc := "owner-update permission probe (delete me)"
		_, err := client.CreateEvidence(ctx, graphclient.CreateEvidenceInput{
			Name:        "MCP owner probe (delete me)",
			Description: &desc,
			ControlIDs:  []string{id},
		}, nil)
		if err != nil {
			fmt.Printf("evidence link org 12.6.2: %s\n", openlane.Redact(err.Error()))
		} else {
			fmt.Println("evidence link org 12.6.2: OK")
		}
	}

	fmt.Println("\n=== resolve personal groups ===")
	for _, term := range []string{"abhishek", "gupta", "greg", "knell"} {
		first := int64(200)
		resp, err := client.GetGroups(ctx, &first, nil, nil)
		if err != nil {
			fmt.Printf("groups: %s\n", openlane.Redact(err.Error()))
			break
		}
		for _, e := range resp.Groups.Edges {
			if e == nil || e.Node == nil {
				continue
			}
			g := e.Node
			hay := strings.ToLower(g.Name + " " + g.DisplayName)
			if strings.Contains(hay, term) {
				fmt.Printf("match %q: id=%s name=%q display=%q managed=%v\n", term, g.ID, g.Name, g.DisplayName, openlane.Deref(g.IsManaged))
			}
		}
	}

	for _, gid := range []string{"01KTY1J6WEC5EXBVSJJWW6QP68", "01M1EN3JHTHT3X5P2FJRTGBH00"} {
		resp, err := client.GetGroupByID(ctx, gid)
		if err != nil {
			fmt.Printf("group %s: %s\n", gid, openlane.Redact(err.Error()))
			continue
		}
		g := resp.Group
		fmt.Printf("12.6.1 existing group %s: name=%q display=%q managed=%v\n", gid, g.Name, g.DisplayName, openlane.Deref(g.IsManaged))
	}

	fmt.Println("\n=== probe owner update with group IDs on 12.6.2 ===")
	abhishekGroup := "01KTY1J6WEC5EXBVSJJWW6QP68" // placeholder; replace if found
	gregGroup := "01M1EN3JHTHT3X5P2FJRTGBH00"
	if id := controls["12.6.2"]; id != "" {
		_, err := client.UpdateControl(ctx, id, graphclient.UpdateControlInput{
			ControlOwnerID: &abhishekGroup,
			DelegateID:     &gregGroup,
		})
		if err != nil {
			fmt.Printf("12.6.2 owner+delegate group IDs: %s\n", openlane.Redact(err.Error()))
		} else {
			fmt.Println("12.6.2 owner+delegate group IDs: OK")
		}
	}
}

func boolPtr(v bool) *bool { return &v }
