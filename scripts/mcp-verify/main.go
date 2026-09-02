package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	root, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, "cwd: %v\n", err)
		os.Exit(1)
	}

	cmd := exec.Command(root+"/bin/openlane-mcp", "serve")
	cmd.Env = os.Environ()
	cmd.Dir = root

	client := mcp.NewClient(&mcp.Implementation{Name: "mcp-verify", Version: "1.0"}, nil)
	session, err := client.Connect(ctx, &mcp.CommandTransport{Command: cmd}, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "connect failed: %v\n", err)
		os.Exit(1)
	}
	defer session.Close()

	tools, err := session.ListTools(ctx, nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "tools/list failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("ok: %d tools registered\n", len(tools.Tools))
	if len(tools.Tools) == 0 {
		os.Exit(1)
	}
}
