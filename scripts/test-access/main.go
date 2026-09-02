//go:build ignore

package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func main() {
	cfg, err := config.FromEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "config: %v\n", err)
		os.Exit(1)
	}
	client, err := openlane.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "client: %v\n", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	first := int64(1)
	resp, err := client.GetControls(ctx, &first, nil, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "GetControls failed: %v\n", openlane.Redact(err.Error()))
		os.Exit(1)
	}

	fmt.Printf("ok: total_count=%d returned=%d\n", resp.Controls.TotalCount, len(resp.Controls.Edges))
	if len(resp.Controls.Edges) > 0 && resp.Controls.Edges[0].Node != nil {
		n := resp.Controls.Edges[0].Node
		title := openlane.Deref(n.Title)
		fmt.Printf("sample: id=%s ref=%s title=%q\n", n.ID, n.RefCode, title)
	}
}
