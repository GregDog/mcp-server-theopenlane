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
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	controlID := "01KZ4X6JF4RSX3NE52W9PG692Z"  // Security Awareness Training
	evidenceID := "01M35G7H965HY5KTRX51PACEJZ" // existing KnowBe4 test evidence

	before := mustGet(client, ctx, evidenceID)
	fmt.Printf("before id=%s controls=%v\n", evidenceID, controlIDs(before))

	upd, updateErr := client.UpdateEvidence(ctx, evidenceID, graphclient.UpdateEvidenceInput{
		AddControlIDs: []string{controlID},
	}, nil)
	if updateErr != nil {
		fmt.Printf("update_error: %s\n", openlane.Redact(updateErr.Error()))
	} else {
		fmt.Printf("update_ok mutation_id=%s\n", upd.UpdateEvidence.Evidence.ID)
	}

	after := mustGet(client, ctx, evidenceID)
	linked := controlIDs(after)
	fmt.Printf("after_get controls=%v linked=%v\n", linked, contains(linked, controlID))

	first := int64(20)
	listed, listErr := client.GetEvidences(ctx, &first, nil, &graphclient.EvidenceWhereInput{
		HasControlsWith: []*graphclient.ControlWhereInput{{ID: &controlID}},
	})
	if listErr != nil {
		fmt.Printf("list_by_control_error: %s\n", openlane.Redact(listErr.Error()))
	} else {
		found := false
		for _, edge := range listed.Evidences.Edges {
			if edge != nil && edge.Node != nil && edge.Node.ID == evidenceID {
				found = true
				break
			}
		}
		fmt.Printf("list_by_control total=%d contains_evidence=%v\n", listed.Evidences.TotalCount, found)
	}

	switch {
	case contains(linked, controlID) && updateErr != nil:
		fmt.Println("RESULT: control linked on get despite mutation error")
	case contains(linked, controlID):
		fmt.Println("RESULT: control linked on get")
	case updateErr != nil:
		fmt.Println("RESULT: mutation error and control not linked on get/list")
	default:
		fmt.Println("RESULT: mutation ok but control not linked on get/list")
	}
}

func mustGet(client openlane.GraphAPI, ctx context.Context, id string) graphclient.GetEvidenceByID_Evidence {
	got, err := client.GetEvidenceByID(ctx, id)
	if err != nil {
		fmt.Printf("get_error id=%s: %s\n", id, openlane.Redact(err.Error()))
		os.Exit(1)
	}
	return got.Evidence
}

func controlIDs(e graphclient.GetEvidenceByID_Evidence) []string {
	var ids []string
	for _, edge := range e.Controls.Edges {
		if edge != nil && edge.Node != nil {
			ids = append(ids, edge.Node.ID)
		}
	}
	return ids
}

func contains(xs []string, want string) bool {
	for _, x := range xs {
		if x == want {
			return true
		}
	}
	return false
}
