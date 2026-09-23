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
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	token := cfg.APIToken
	prefix := token
	if len(prefix) > 4 {
		prefix = prefix[:4]
	}
	fmt.Printf("token_prefix=%s… org=%s\n", prefix, cfg.OrganizationID)

	taskResp, taskErr := client.CreateTask(ctx, graphclient.CreateTaskInput{
		Title: "MCP write probe (delete me)",
		Tags:  []string{"mcp-test"},
	})
	if taskErr != nil {
		fmt.Println("createTask:", openlane.Redact(taskErr.Error()))
	} else {
		fmt.Println("createTask_ok:", taskResp.CreateTask.Task.ID)
	}

	desc := "Live token probe. Safe to delete."
	cases := []struct {
		label string
		input graphclient.CreateEvidenceInput
	}{
		{"create_minimal", graphclient.CreateEvidenceInput{Name: "MCP test: KnowBe4 training (delete me)", Description: &desc}},
		{"create_with_source", graphclient.CreateEvidenceInput{Name: "MCP test: KnowBe4 training (delete me)", Description: &desc, Source: strPtr("KnowBe4"), Tags: []string{"mcp-test"}}},
		{"create_with_control", graphclient.CreateEvidenceInput{
			Name: "MCP test: KnowBe4 training (delete me)", Description: &desc, Source: strPtr("KnowBe4"),
			Tags: []string{"mcp-test"}, ControlIDs: []string{"01KZ4X6JF4RSX3NE52W9PG692Z"},
		}},
	}
	var evidenceID string
	for _, c := range cases {
		resp, err := client.CreateEvidence(ctx, c.input, nil)
		if err != nil {
			fmt.Printf("%s: %s\n", c.label, openlane.Redact(err.Error()))
			continue
		}
		evidenceID = resp.CreateEvidence.Evidence.ID
		fmt.Printf("%s_ok: %s\n", c.label, evidenceID)
	}

	if evidenceID == "" {
		return
	}
	controlID := "01KZ4X6PQ969A9HASAJGERB78A" // OL-12.06 Forecasting from list probe
	upd, err := client.UpdateEvidence(ctx, evidenceID, graphclient.UpdateEvidenceInput{
		AddControlIDs: []string{controlID},
	}, nil)
	if err != nil {
		fmt.Printf("update_add_control_ids: %s\n", openlane.Redact(err.Error()))
		return
	}
	fmt.Printf("update_add_control_ids_ok: %s\n", upd.UpdateEvidence.Evidence.ID)

	got, err := client.GetEvidenceByID(ctx, evidenceID)
	if err != nil {
		fmt.Printf("get_evidence: %s\n", openlane.Redact(err.Error()))
		return
	}
	var linked []string
	for _, edge := range got.Evidence.Controls.Edges {
		if edge != nil && edge.Node != nil {
			linked = append(linked, edge.Node.ID)
		}
	}
	fmt.Printf("get_controls=%v\n", linked)
}

func strPtr(s string) *string { return &s }
