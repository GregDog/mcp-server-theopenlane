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
	ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
	defer cancel()

	knownIDs := []string{
		"01M35G72VM2ZV2YHCM75DE1J6Q",
		"01M35G7343ZQ6WHGRWYCPMCR6C",
		"01M35G7H22KTBQ1T4DGWPYWHAJ",
		"01M35G7H965HY5KTRX51PACEJZ",
		"01M35G7GRFPEBEJW7EQCM5KDEJ", // task, not evidence - will fail gracefully
	}
	fmt.Println("=== get by known IDs from test output ===")
	for _, id := range knownIDs {
		got, err := client.GetEvidenceByID(ctx, id)
		if err != nil {
			fmt.Printf("id=%s missing_or_error: %s\n", id, openlane.Redact(err.Error()))
			continue
		}
		printDetail(got.Evidence)
	}

	names := []string{
		"MCP test: KnowBe4 training (delete me)",
		"MCP write probe (delete me)",
		"MCP test: link probe (delete me)",
	}
	fmt.Println("\n=== search by name ===")
	for _, name := range names {
		first := int64(20)
		where := &graphclient.EvidenceWhereInput{NameContainsFold: &name}
		resp, err := client.GetEvidences(ctx, &first, nil, where)
		if err != nil {
			fmt.Printf("name=%q error: %s\n", name, openlane.Redact(err.Error()))
			continue
		}
		fmt.Printf("name=%q matches=%d total=%d\n", name, len(resp.Evidences.Edges), resp.Evidences.TotalCount)
		for _, edge := range resp.Evidences.Edges {
			if edge == nil || edge.Node == nil {
				continue
			}
			fmt.Printf("  list id=%s name=%q source=%q\n", edge.Node.ID, edge.Node.Name, openlane.Deref(edge.Node.Source))
		}
	}

	fmt.Println("\n=== recent evidence containing 'MCP' ===")
	first := int64(50)
	all, err := client.GetEvidences(ctx, &first, nil, nil)
	if err != nil {
		fmt.Println("list_all:", openlane.Redact(err.Error()))
		return
	}
	fmt.Printf("total=%d returned=%d\n", all.Evidences.TotalCount, len(all.Evidences.Edges))
	for _, edge := range all.Evidences.Edges {
		if edge == nil || edge.Node == nil {
			continue
		}
		if strings.Contains(strings.ToLower(edge.Node.Name), "mcp") {
			fmt.Printf("  list id=%s name=%q source=%q\n", edge.Node.ID, edge.Node.Name, openlane.Deref(edge.Node.Source))
		}
	}
}

func printDetail(e graphclient.GetEvidenceByID_Evidence) {
	var controls []string
	for _, c := range e.Controls.Edges {
		if c != nil && c.Node != nil {
			controls = append(controls, c.Node.ID)
		}
	}
	fmt.Printf("id=%s name=%q source=%q controls=%v status=%s\n",
		e.ID, e.Name, openlane.Deref(e.Source), controls, openlane.Format(e.Status))
}
