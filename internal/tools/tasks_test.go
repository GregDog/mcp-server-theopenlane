package tools

import (
	"context"
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

func TestListTasks(t *testing.T) {
	api := &fakeAPI{
		tasks: &graphclient.GetTasks{
			Tasks: graphclient.GetTasks_Tasks{
				TotalCount: 1,
				PageInfo: graphclient.GetTasks_Tasks_PageInfo{
					HasNextPage: false,
				},
				Edges: []*graphclient.GetTasks_Tasks_Edges{
					{
						Node: &graphclient.GetTasks_Tasks_Edges_Node{
							ID:        "task_1",
							DisplayID: "TSK-1",
							Title:     "Review access",
							Status:    "OPEN",
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api}
	_, page, err := h.listTasks(context.Background(), nil, listInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 || page.Items[0].ID != "task_1" {
		t.Fatalf("unexpected page: %+v", page)
	}
}

func TestGetTaskRequiresID(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.getTask(context.Background(), nil, getInput{})
	if err != errIDRequired {
		t.Fatalf("expected errIDRequired, got %v", err)
	}
}
