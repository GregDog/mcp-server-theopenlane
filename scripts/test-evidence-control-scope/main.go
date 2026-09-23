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
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	first := int64(10)
	resp, err := client.GetControls(ctx, &first, nil, nil)
	if err != nil {
		fmt.Println("list_controls:", openlane.Redact(err.Error()))
		os.Exit(1)
	}

	desc := "control link probe"
	for _, edge := range resp.Controls.Edges {
		if edge == nil || edge.Node == nil {
			continue
		}
		c := edge.Node
		title := openlane.Deref(c.Title)
		fmt.Printf("try_control id=%s ref=%s title=%q framework=%q source=%s\n",
			c.ID, c.RefCode, title, openlane.Deref(c.ReferenceFramework), openlane.Format(c.Source))

		_, err := client.CreateEvidence(ctx, graphclient.CreateEvidenceInput{
			Name:        "MCP control link probe (delete me)",
			Description: &desc,
			ControlIDs:  []string{c.ID},
		}, nil)
		if err != nil {
			msg := openlane.Redact(err.Error())
			if strings.Contains(msg, "UNAUTHORIZED") {
				fmt.Println("  result=UNAUTHORIZED")
			} else {
				fmt.Printf("  result=%s\n", msg)
			}
			continue
		}
		fmt.Println("  result=OK")
		break
	}
}
