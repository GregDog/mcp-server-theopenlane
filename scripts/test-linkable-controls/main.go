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

	ownerSet := true
	q := "12.6.2"
	first := int64(20)
	resp, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
		OwnerIDNotNil:       &ownerSet,
		RefCodeContainsFold: &q,
	})
	if err != nil {
		fmt.Println("search:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("linkable_12_6_2_count=%d\n", len(resp.Controls.Edges))
	for _, e := range resp.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		fmt.Printf("id=%s ref=%s source=%s title=%q\n", n.ID, n.RefCode, openlane.Format(n.Source), openlane.Deref(n.Title))
	}

	resp2, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
		OwnerIDNotNil: &ownerSet,
	})
	if err != nil {
		fmt.Println("list:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("total_linkable_sample=%d\n", len(resp2.Controls.Edges))
}
