package tools

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type programItem struct {
	ID            string `json:"id"`
	DisplayID     string `json:"display_id,omitempty"`
	Name          string `json:"name,omitempty"`
	Status        string `json:"status,omitempty"`
	ProgramKind   string `json:"program_kind,omitempty"`
	FrameworkName string `json:"framework_name,omitempty"`
	Description   string `json:"description,omitempty"`
	StartDate     string `json:"start_date,omitempty"`
	EndDate       string `json:"end_date,omitempty"`
	AuditorReady  *bool  `json:"auditor_ready,omitempty"`
	AuditFirm     string `json:"audit_firm,omitempty"`
}

func registerPrograms(server *mcp.Server, h *handlers) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_programs_list",
		Title:       "List Openlane programs",
		Description: "List compliance programs in the configured Openlane organization. Results are paginated.",
		Annotations: readOnly(),
	}, h.listPrograms)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "openlane_program_get",
		Title:       "Get an Openlane program",
		Description: "Get a single program by ID.",
		Annotations: readOnly(),
	}, h.getProgram)
}

func (h *handlers) listPrograms(ctx context.Context, _ *mcp.CallToolRequest, in listInput) (*mcp.CallToolResult, openlane.Page[programItem], error) {
	first, after := pageArgs(in.Limit, in.Cursor)
	resp, err := h.api.GetPrograms(ctx, &first, after, nil)
	if err != nil {
		return nil, openlane.Page[programItem]{}, openlane.APIError(err)
	}
	items := make([]programItem, 0, len(resp.Programs.Edges))
	for _, e := range resp.Programs.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, programItem{
			ID:            n.ID,
			DisplayID:     n.DisplayID,
			Name:          n.Name,
			Status:        openlane.Format(n.Status),
			ProgramKind:   openlane.Deref(n.ProgramKindName),
			FrameworkName: openlane.Deref(n.FrameworkName),
		})
	}
	return nil, openlane.Page[programItem]{
		Items:      items,
		NextCursor: resp.Programs.PageInfo.EndCursor,
		HasMore:    resp.Programs.PageInfo.HasNextPage,
		TotalCount: resp.Programs.TotalCount,
	}, nil
}

func (h *handlers) getProgram(ctx context.Context, _ *mcp.CallToolRequest, in getInput) (*mcp.CallToolResult, programItem, error) {
	if in.ID == "" {
		return nil, programItem{}, errIDRequired
	}
	resp, err := h.api.GetProgramByID(ctx, in.ID)
	if err != nil {
		return nil, programItem{}, openlane.APIError(err)
	}
	p := resp.Program
	ready := p.AuditorReady
	return nil, programItem{
		ID:            p.ID,
		DisplayID:     p.DisplayID,
		Name:          p.Name,
		Status:        openlane.Format(p.Status),
		ProgramKind:   openlane.Deref(p.ProgramKindName),
		FrameworkName: openlane.Deref(p.FrameworkName),
		Description:   openlane.Deref(p.Description),
		StartDate:     openlane.Format(p.StartDate),
		EndDate:       openlane.Format(p.EndDate),
		AuditorReady:  &ready,
		AuditFirm:     openlane.Deref(p.AuditFirm),
	}, nil
}
