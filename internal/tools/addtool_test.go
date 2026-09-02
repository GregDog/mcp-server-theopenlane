package tools

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func TestAddToolOmitsOutputSchema(t *testing.T) {
	server := mcp.NewServer(&mcp.Implementation{Name: "test", Version: "1"}, nil)
	addTool(server, &mcp.Tool{Name: "echo"}, func(_ context.Context, _ *mcp.CallToolRequest, in struct {
		Message string `json:"message"`
	}) (*mcp.CallToolResult, struct {
		Message string `json:"message"`
	}, error) {
		return nil, struct {
			Message string `json:"message"`
		}{Message: in.Message}, nil
	})

	t1, t2 := mcp.NewInMemoryTransports()
	ctx := context.Background()
	ss, err := server.Connect(ctx, t1, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer ss.Close()
	cs, err := mcp.NewClient(&mcp.Implementation{Name: "client", Version: "1"}, nil).Connect(ctx, t2, nil)
	if err != nil {
		t.Fatal(err)
	}
	defer cs.Close()

	tools, err := cs.ListTools(ctx, nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(tools.Tools) != 1 {
		t.Fatalf("tools: %d", len(tools.Tools))
	}
	if tools.Tools[0].OutputSchema != nil {
		t.Fatal("expected output schema to be omitted")
	}
	b, err := json.Marshal(tools)
	if err != nil {
		t.Fatal(err)
	}
	if len(b) > 500 {
		t.Fatalf("unexpected tools/list payload size: %d", len(b))
	}
}
