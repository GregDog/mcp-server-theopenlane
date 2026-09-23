package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
	"github.com/theopenlane/core/common/enums"
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
	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	name := "PCI DSS"
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	first := int64(20)
	programs, err := client.GetPrograms(ctx, &first, nil, &graphclient.ProgramWhereInput{
		NameContainsFold: &name,
	})
	if err != nil {
		fmt.Println("programs:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("programs_matching_%q=%d\n", name, len(programs.Programs.Edges))
	for _, e := range programs.Programs.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		p := e.Node
		fmt.Printf("program id=%s name=%q framework=%q\n", p.ID, p.Name, openlane.Deref(p.FrameworkName))
	}

	first = int64(50)
	all, err := client.GetControls(ctx, &first, nil, nil)
	if err != nil {
		fmt.Println("controls:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	bySource := map[string]int{}
	for _, e := range all.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		bySource[openlane.Format(e.Node.Source)]++
	}
	fmt.Printf("control_sample_total=%d by_source=%v\n", len(all.Controls.Edges), bySource)

	framework := enums.ControlSourceFramework
	linkable, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
		SourceNotIn: []enums.ControlSource{framework},
	})
	if err != nil {
		fmt.Println("linkable:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("linkable_controls_sample=%d\n", len(linkable.Controls.Edges))
	for _, e := range linkable.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		fmt.Printf("linkable id=%s ref=%s source=%s title=%q\n", n.ID, n.RefCode, openlane.Format(n.Source), openlane.Deref(n.Title))
	}

	if len(programs.Programs.Edges) == 0 || programs.Programs.Edges[0] == nil || programs.Programs.Edges[0].Node == nil {
		return
	}
	programID := programs.Programs.Edges[0].Node.ID
	inProgram, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
		HasProgramsWith: []*graphclient.ProgramWhereInput{{ID: &programID}},
	})
	if err != nil {
		fmt.Println("program_controls:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("controls_in_program=%d total=%d\n", len(inProgram.Controls.Edges), inProgram.Controls.TotalCount)
	for _, e := range inProgram.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		fmt.Printf("program_control id=%s ref=%s source=%s owner=%s display=%s title=%q\n",
			n.ID, n.RefCode, openlane.Format(n.Source), openlane.Deref(n.OwnerID), n.DisplayID, openlane.Deref(n.Title))
	}

	ref := "12.6.2"
	byRef, err := client.GetControls(ctx, &first, nil, &graphclient.ControlWhereInput{
		RefCodeContainsFold: &ref,
	})
	if err != nil {
		fmt.Println("ref_search:", openlane.Redact(err.Error()))
		os.Exit(1)
	}
	fmt.Printf("controls_matching_ref_%s=%d total=%d\n", ref, len(byRef.Controls.Edges), byRef.Controls.TotalCount)
	for _, e := range byRef.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		fmt.Printf("ref_control id=%s ref=%s source=%s owner=%s display=%s title=%q\n",
			n.ID, n.RefCode, openlane.Format(n.Source), openlane.Deref(n.OwnerID), n.DisplayID, openlane.Deref(n.Title))
	}
}
