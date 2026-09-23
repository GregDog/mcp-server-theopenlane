package main

import (
	"context"
	"encoding/json"
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

	controlID := "01KZ4X6JF4RSX3NE52W9PG692Z"
	desc := "relationship write probe"

	// 1) evidence + controlIDs
	in := graphclient.CreateEvidenceInput{
		Name: "MCP rel probe (delete me)", Description: &desc, ControlIDs: []string{controlID},
	}
	b, _ := json.Marshal(in)
	fmt.Printf("evidence_control_payload=%s\n", b)
	try("create_evidence_controlIDs", func() error {
		_, err := client.CreateEvidence(ctx, in, nil)
		return err
	})

	// 2) risk + entityIDs (known working pattern)
	first := int64(1)
	ents, err := client.GetEntities(ctx, &first, nil, nil)
	if err == nil && len(ents.Entities.Edges) > 0 && ents.Entities.Edges[0].Node != nil {
		entID := ents.Entities.Edges[0].Node.ID
		try("create_risk_entityIDs", func() error {
			_, err := client.CreateRisk(ctx, graphclient.CreateRiskInput{
				Name: "MCP rel probe risk (delete me)", EntityIDs: []string{entID},
			})
			return err
		})
	}

	// 3) control metadata write
	note := "MCP control metadata probe"
	try("update_control_description", func() error {
		_, err := client.UpdateControl(ctx, controlID, graphclient.UpdateControlInput{Description: &note})
		return err
	})

	// 4) program link if any program exists
	progs, err := client.GetPrograms(ctx, &first, nil, nil)
	if err == nil && len(progs.Programs.Edges) > 0 && progs.Programs.Edges[0].Node != nil {
		progID := progs.Programs.Edges[0].Node.ID
		try("create_evidence_programIDs", func() error {
			_, err := client.CreateEvidence(ctx, graphclient.CreateEvidenceInput{
				Name: "MCP rel probe program (delete me)", Description: &desc, ProgramIDs: []string{progID},
			}, nil)
			return err
		})
	}
}

func try(label string, fn func() error) {
	if err := fn(); err != nil {
		fmt.Printf("%s: %s\n", label, openlane.Redact(err.Error()))
		return
	}
	fmt.Printf("%s: ok\n", label)
}
