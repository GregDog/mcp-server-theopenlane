package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func main() {
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
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	fromID := envOr("FROM_CONTROL_ID", "01KTY1WZBS3V5C6ZY2YJ3ANMED")
	toID := envOr("TO_CONTROL_ID", "01KTY1WZBZZ17NK5C7CJXGJW97")
	mt := enums.MappingTypeIntersect
	src := enums.MappingSourceManual
	relation := "MCP mapped-control probe"

	created, err := client.CreateMappedControl(ctx, graphclient.CreateMappedControlInput{
		MappingType:    &mt,
		Source:         &src,
		Relation:       &relation,
		FromControlIDs: []string{fromID},
		ToControlIDs:   []string{toID},
	})
	if err != nil {
		fmt.Println("create:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	mapID := created.CreateMappedControl.MappedControl.ID
	fmt.Printf("created=%s\n", mapID)

	got, err := client.GetMappedControlByID(ctx, mapID)
	if err != nil {
		fmt.Println("get:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("get_from=%d get_to=%d\n", len(got.MappedControl.FromControls.Edges), len(got.MappedControl.ToControls.Edges))

	first := int64(10)
	cw := []*graphclient.ControlWhereInput{{ID: &fromID}}
	list, err := client.GetMappedControls(ctx, &first, nil, &graphclient.MappedControlWhereInput{
		Or: []*graphclient.MappedControlWhereInput{
			{HasFromControlsWith: cw},
			{HasToControlsWith: cw},
		},
	})
	if err != nil {
		fmt.Println("list:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("list_for_control=%d\n", list.MappedControls.TotalCount)

	deletedID, err := client.DeleteMappedControl(ctx, mapID)
	if err != nil {
		fmt.Println("delete:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("deleted=%s\n", deletedID)
	fmt.Println("RESULT: ok")
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
