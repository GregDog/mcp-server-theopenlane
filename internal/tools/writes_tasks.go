package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type taskItem struct {
	ID        string `json:"id"`
	DisplayID string `json:"display_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Status    string `json:"status,omitempty"`
	Details   string `json:"details,omitempty"`
	Due       string `json:"due,omitempty"`
}

type createTaskInput struct {
	Title   string   `json:"title" jsonschema:"Task title."`
	Details string   `json:"details,omitempty" jsonschema:"Task details."`
	Status  string   `json:"status,omitempty" jsonschema:"Task status enum value."`
	Tags    []string `json:"tags,omitempty" jsonschema:"Tags to apply."`
}

type updateTaskInput struct {
	ID      string   `json:"id" jsonschema:"Task ID to update."`
	Title   string   `json:"title,omitempty" jsonschema:"Updated title."`
	Details string   `json:"details,omitempty" jsonschema:"Updated details."`
	Status  string   `json:"status,omitempty" jsonschema:"Updated task status enum value."`
	Tags    []string `json:"tags,omitempty" jsonschema:"Replace tags with this list."`
}

func registerWriteTasks(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_task_create",
		Title:       "Create an Openlane task",
		Description: "Create a task in Openlane. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.createTask)

	addTool(server, &mcp.Tool{
		Name:        "openlane_task_update",
		Title:       "Update an Openlane task",
		Description: "Update a task by ID. Requires write mode.",
		Annotations: writeAnnotations(),
	}, h.updateTask)
}

func (h *handlers) createTask(ctx context.Context, _ *mcp.CallToolRequest, in createTaskInput) (*mcp.CallToolResult, taskItem, error) {
	if in.Title == "" {
		return nil, taskItem{}, errTitleRequired
	}
	input := graphclient.CreateTaskInput{
		Title: in.Title,
		Tags:  in.Tags,
	}
	if in.Details != "" {
		input.Details = &in.Details
	}
	if in.Status != "" {
		input.Status = taskStatus(in.Status)
	}

	resp, err := h.api.CreateTask(ctx, input)
	if err != nil {
		return nil, taskItem{}, openlane.APIError(err)
	}
	return nil, mapCreatedTask(resp.CreateTask.Task), nil
}

func (h *handlers) updateTask(ctx context.Context, _ *mcp.CallToolRequest, in updateTaskInput) (*mcp.CallToolResult, taskItem, error) {
	if in.ID == "" {
		return nil, taskItem{}, errIDRequired
	}
	input := graphclient.UpdateTaskInput{}
	if in.Title != "" {
		input.Title = &in.Title
	}
	if in.Details != "" {
		input.Details = &in.Details
	}
	if in.Status != "" {
		input.Status = taskStatus(in.Status)
	}
	if len(in.Tags) > 0 {
		input.Tags = in.Tags
	}
	if isEmptyUpdateTask(input) {
		return nil, taskItem{}, errUpdateFieldsRequired
	}

	resp, err := h.api.UpdateTask(ctx, in.ID, input)
	if err != nil {
		return nil, taskItem{}, openlane.APIError(err)
	}
	return nil, mapUpdatedTask(resp.UpdateTask.Task), nil
}

func mapCreatedTask(t graphclient.CreateTask_CreateTask_Task) taskItem {
	return taskItem{
		ID:        t.ID,
		DisplayID: t.DisplayID,
		Title:     t.Title,
		Status:    openlane.Format(t.Status),
		Details:   openlane.Deref(t.Details),
		Due:       openlane.Format(t.Due),
	}
}

func mapUpdatedTask(t graphclient.UpdateTask_UpdateTask_Task) taskItem {
	return taskItem{
		ID:        t.ID,
		DisplayID: t.DisplayID,
		Title:     t.Title,
		Status:    openlane.Format(t.Status),
		Details:   openlane.Deref(t.Details),
		Due:       openlane.Format(t.Due),
	}
}

func isEmptyUpdateTask(in graphclient.UpdateTaskInput) bool {
	return in.Title == nil &&
		in.Details == nil &&
		in.Status == nil &&
		len(in.Tags) == 0
}
