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
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	client, err := openlane.New(cfg)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	refs := []string{"12.6.1", "12.6.2", "12.6.3", "12.6.3.1", "12.6.3.2"}
	ownerSet := true
	for _, ref := range refs {
		first := int64(5)
		q := ref
		resp, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
			OwnerIDNotNil:       &ownerSet,
			RefCodeContainsFold: &q,
		})
		if err != nil {
			fmt.Printf("%s error: %v\n", ref, err)
			continue
		}
		for _, e := range resp.Controls.Edges {
			if e == nil || e.Node == nil {
				continue
			}
			n := e.Node
			if !strings.EqualFold(n.RefCode, ref) {
				continue
			}
			cid := n.ID
			fmt.Printf("\n=== %s %s (%s) ===\n", n.RefCode, n.DisplayID, cid)
			fmt.Printf("title=%q\n", openlane.Deref(n.Title))

			cw := []*graphclient.ControlWhereInput{{ID: &cid}}
			pfirst := int64(10)
			progs, _ := client.GetPrograms(ctx, &pfirst, nil, &graphclient.ProgramWhereInput{HasControlsWith: cw})
			if progs != nil {
				fmt.Printf("programs=%d", progs.Programs.TotalCount)
				for _, pe := range progs.Programs.Edges {
					if pe != nil && pe.Node != nil {
						fmt.Printf(" [%s]", pe.Node.Name)
					}
				}
				fmt.Println()
			}

			efirst := int64(20)
			ev, _ := client.GetEvidences(ctx, &efirst, nil, &graphclient.EvidenceWhereInput{HasControlsWith: cw})
			if ev != nil {
				fmt.Printf("evidence=%d\n", ev.Evidences.TotalCount)
				for _, ee := range ev.Evidences.Edges {
					if ee == nil || ee.Node == nil {
						continue
					}
					fmt.Printf("  - %s (%s) source=%s\n", ee.Node.Name, ee.Node.ID, openlane.Deref(ee.Node.Source))
				}
			}

			ifirst := int64(10)
			impl, _ := client.GetControlImplementations(ctx, &ifirst, nil, &graphclient.ControlImplementationWhereInput{HasControlsWith: cw})
			if impl != nil {
				fmt.Printf("implementations=%d\n", impl.ControlImplementations.TotalCount)
			}
		}
	}
}
