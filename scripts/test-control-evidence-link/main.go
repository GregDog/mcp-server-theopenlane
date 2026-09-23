package main

import (
	"context"
	"fmt"
	"os"
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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	controlID := "01KZ4X6JF4RSX3NE52W9PG692Z"
	evidenceID := "01M35G7H965HY5KTRX51PACEJZ"

	ctl, err := client.GetControlByID(ctx, controlID)
	if err != nil {
		fmt.Println("get_control:", openlane.Redact(err.Error()))
	} else {
		n := ctl.Control
		fmt.Printf("control title=%q ref=%s framework=%q source=%s standard=%s\n",
			openlane.Deref(n.Title), n.RefCode, openlane.Deref(n.ReferenceFramework), openlane.Format(n.Source), openlane.Deref(n.StandardID))
	}

	desc := "MCP metadata probe"
	_, err = client.UpdateControl(ctx, controlID, graphclient.UpdateControlInput{Description: &desc})
	if err != nil {
		fmt.Println("control_description_update:", openlane.Redact(err.Error()))
	} else {
		fmt.Println("control_description_update: ok")
	}

	fmt.Println("=== updateEvidence addControlIDs ===")
	_, evErr := client.UpdateEvidence(ctx, evidenceID, graphclient.UpdateEvidenceInput{
		AddControlIDs: []string{controlID},
	}, nil)
	if evErr != nil {
		fmt.Println("evidence_side:", openlane.Redact(evErr.Error()))
	} else {
		fmt.Println("evidence_side: ok")
	}
	printLink(client, ctx, evidenceID, controlID)

	fmt.Println("=== updateControl addEvidenceIDs ===")
	_, ctlErr := client.UpdateControl(ctx, controlID, graphclient.UpdateControlInput{
		AddEvidenceIDs: []string{evidenceID},
	})
	if ctlErr != nil {
		fmt.Println("control_side:", openlane.Redact(ctlErr.Error()))
	} else {
		fmt.Println("control_side: ok")
	}
	printLink(client, ctx, evidenceID, controlID)
}

func printLink(client openlane.GraphAPI, ctx context.Context, evidenceID, controlID string) {
	got, err := client.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		fmt.Println("verify:", openlane.Redact(err.Error()))
		return
	}
	var linked []string
	for _, edge := range got.Evidence.Controls.Edges {
		if edge != nil && edge.Node != nil {
			linked = append(linked, edge.Node.ID)
		}
	}
	fmt.Printf("verify controls=%v linked=%v\n", linked, contains(linked, controlID))
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
