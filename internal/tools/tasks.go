package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func registerTasks(server *mcp.Server, h *handlers) {
	addTool(server, &mcp.Tool{
		Name:        "openlane_tasks_list",
		Title:       "List Openlane tasks",
		Description: "List tasks in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listTasks)

	addTool(server, &mcp.Tool{
		Name:        "openlane_task_get",
		Title:       "Get an Openlane task",
		Description: "Get a single task by ID.",
		Annotations: readOnly(),
	}, h.getTask)
}

func (h *handlers) listTasks(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[taskItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetTasks(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[taskItem]{}, openlane.APIError(err)
	}
	items := make([]taskItem, 0, len(resp.Tasks.Edges))
	for _, e := range resp.Tasks.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		items = append(items, mapListTask(e.Node))
	}
	return nil, openlane.Page[taskItem]{
		Items:      items,
		NextCursor: resp.Tasks.PageInfo.EndCursor,
		HasMore:    resp.Tasks.PageInfo.HasNextPage,
		TotalCount: resp.Tasks.TotalCount,
	}, nil
}

func (h *handlers) getTask(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, taskItem, error) {
	if in.ID == "" {
		return nil, taskItem{}, errIDRequired
	}
	resp, err := h.api.GetTaskByID(ctx, in.ID)
	if err != nil {
		return nil, taskItem{}, openlane.APIError(err)
	}
	return nil, mapGetTask(resp.Task), nil
}

func mapListTask(n *graphclient.GetTasks_Tasks_Edges_Node) taskItem {
	return taskItem{
		ID:        n.ID,
		DisplayID: n.DisplayID,
		Title:     n.Title,
		Status:    openlane.Format(n.Status),
		Details:   openlane.Deref(n.Details),
		Due:       openlane.Format(n.Due),
	}
}

func mapGetTask(t graphclient.GetTaskByID_Task) taskItem {
	return taskItem{
		ID:        t.ID,
		DisplayID: t.DisplayID,
		Title:     t.Title,
		Status:    openlane.Format(t.Status),
		Details:   openlane.Deref(t.Details),
		Due:       openlane.Format(t.Due),
	}
}
