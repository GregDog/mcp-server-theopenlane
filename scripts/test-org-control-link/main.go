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

	orgControl := "01KTY1WZBZZ17NK5C7CJXGJW97" // 12.6.2 org-owned in PCI program
	catControl := "01K39TXRZQZBZS3S1YMXJVE914" // 12.6.2 catalog/system
	desc := "probe"
	for name, id := range map[string]string{"org_owned_12_6_2": orgControl, "catalog_12_6_2": catControl} {
		resp, err := client.CreateEvidence(ctx, graphclient.CreateEvidenceInput{
			Name:        "MCP link probe (delete me)",
			Description: &desc,
			ControlIDs:  []string{id},
		}, nil)
		if err != nil {
			fmt.Printf("%s: %s\n", name, openlane.Redact(err.Error()))
			continue
		}
		fmt.Printf("%s: OK evidence_id=%s\n", name, resp.CreateEvidence.Evidence.ID)
	}
}
