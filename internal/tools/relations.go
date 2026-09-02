package tools

import (
	"context"
	"sync"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

const relSummaryLimit int64 = 8

// relSummary is a compact bounded view of a related object collection.
type relSummary[T any] struct {
	Count int64 `json:"count"`
	Items []T   `json:"items,omitempty"`
}

func relFirst() int64 { return relSummaryLimit }

type idNameRef struct {
	ID     string `json:"id"`
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type findingRef struct {
	ID       string `json:"id"`
	Name     string `json:"name,omitempty"`
	Severity string `json:"severity,omitempty"`
	Open     *bool  `json:"open,omitempty"`
	Status   string `json:"status,omitempty"`
}

type controlRef struct {
	ID      string `json:"id"`
	RefCode string `json:"ref_code,omitempty"`
	Title   string `json:"title,omitempty"`
	Status  string `json:"status,omitempty"`
}

type implementationRef struct {
	ID       string `json:"id"`
	Status   string `json:"status,omitempty"`
	Verified *bool  `json:"verified,omitempty"`
}

type remediationRef struct {
	ID     string `json:"id"`
	Title  string `json:"title,omitempty"`
	Intent string `json:"intent,omitempty"`
	DueAt  string `json:"due_at,omitempty"`
	State  string `json:"state,omitempty"`
}

type taskRef struct {
	ID        string `json:"id"`
	DisplayID string `json:"display_id,omitempty"`
	Title     string `json:"title,omitempty"`
	Status    string `json:"status,omitempty"`
	Due       string `json:"due,omitempty"`
}

type campaignRef struct {
	ID        string `json:"id"`
	DisplayID string `json:"display_id,omitempty"`
	Name      string `json:"name,omitempty"`
	Status    string `json:"status,omitempty"`
}

type assessmentResponseRef struct {
	ID     string `json:"id"`
	Status string `json:"status,omitempty"`
	Email  string `json:"email,omitempty"`
	Due    string `json:"due_date,omitempty"`
}

type vulnerabilityRef struct {
	ID          string `json:"id"`
	DisplayName string `json:"display_name,omitempty"`
	CveID       string `json:"cve_id,omitempty"`
	Severity    string `json:"severity,omitempty"`
}

func runRelJobs(jobs ...func()) {
	var wg sync.WaitGroup
	for _, job := range jobs {
		wg.Add(1)
		go func(j func()) {
			defer wg.Done()
			j()
		}(job)
	}
	wg.Wait()
}

func (h *handlers) fetchEvidences(ctx context.Context, where *graphclient.EvidenceWhereInput) *relSummary[idNameRef] {
	first := relFirst()
	resp, err := h.api.GetEvidences(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]idNameRef, 0, len(resp.Evidences.Edges))
	for _, e := range resp.Evidences.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, idNameRef{ID: n.ID, Name: n.Name, Status: openlane.Format(n.Status)})
	}
	return &relSummary[idNameRef]{Count: resp.Evidences.TotalCount, Items: items}
}

func (h *handlers) fetchFindings(ctx context.Context, where *graphclient.FindingWhereInput) *relSummary[findingRef] {
	first := relFirst()
	resp, err := h.api.GetFindings(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]findingRef, 0, len(resp.Findings.Edges))
	for _, e := range resp.Findings.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, findingRef{
			ID:       n.ID,
			Name:     openlane.Deref(n.DisplayName),
			Severity: openlane.Deref(n.Severity),
			Open:     n.Open,
			Status:   openlane.Deref(n.FindingStatusName),
		})
	}
	return &relSummary[findingRef]{Count: resp.Findings.TotalCount, Items: items}
}

func (h *handlers) fetchRisks(ctx context.Context, where *graphclient.RiskWhereInput) *relSummary[idNameRef] {
	first := relFirst()
	resp, err := h.api.GetRisks(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]idNameRef, 0, len(resp.Risks.Edges))
	for _, e := range resp.Risks.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, idNameRef{ID: n.ID, Name: n.Name, Status: openlane.Format(n.Status)})
	}
	return &relSummary[idNameRef]{Count: resp.Risks.TotalCount, Items: items}
}

func (h *handlers) fetchTasks(ctx context.Context, where *graphclient.TaskWhereInput) *relSummary[taskRef] {
	first := relFirst()
	resp, err := h.api.GetTasks(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]taskRef, 0, len(resp.Tasks.Edges))
	for _, e := range resp.Tasks.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, taskRef{
			ID:        n.ID,
			DisplayID: n.DisplayID,
			Title:     n.Title,
			Status:    openlane.Format(n.Status),
			Due:       openlane.Format(n.Due),
		})
	}
	return &relSummary[taskRef]{Count: resp.Tasks.TotalCount, Items: items}
}

func (h *handlers) fetchRemediations(ctx context.Context, where *graphclient.RemediationWhereInput) *relSummary[remediationRef] {
	first := relFirst()
	resp, err := h.api.GetRemediations(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]remediationRef, 0, len(resp.Remediations.Edges))
	for _, e := range resp.Remediations.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, remediationRef{
			ID:     n.ID,
			Title:  openlane.Deref(n.Title),
			Intent: openlane.Deref(n.Intent),
			DueAt:  openlane.Format(n.DueAt),
			State:  openlane.Deref(n.State),
		})
	}
	return &relSummary[remediationRef]{Count: resp.Remediations.TotalCount, Items: items}
}

func (h *handlers) fetchPrograms(ctx context.Context, where *graphclient.ProgramWhereInput) *relSummary[idNameRef] {
	first := relFirst()
	resp, err := h.api.GetPrograms(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]idNameRef, 0, len(resp.Programs.Edges))
	for _, e := range resp.Programs.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, idNameRef{ID: n.ID, Name: n.Name, Status: openlane.Format(n.Status)})
	}
	return &relSummary[idNameRef]{Count: resp.Programs.TotalCount, Items: items}
}

func (h *handlers) fetchControls(ctx context.Context, where *graphclient.ControlWhereInput) *relSummary[controlRef] {
	first := relFirst()
	resp, err := h.api.GetControls(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]controlRef, 0, len(resp.Controls.Edges))
	for _, e := range resp.Controls.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, controlRef{
			ID:      n.ID,
			RefCode: n.RefCode,
			Title:   openlane.Deref(n.Title),
			Status:  openlane.Format(n.Status),
		})
	}
	return &relSummary[controlRef]{Count: resp.Controls.TotalCount, Items: items}
}

func (h *handlers) fetchEntities(ctx context.Context, where *graphclient.EntityWhereInput) *relSummary[idNameRef] {
	first := relFirst()
	resp, err := h.api.GetEntities(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]idNameRef, 0, len(resp.Entities.Edges))
	for _, e := range resp.Entities.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		name := openlane.Deref(n.Name)
		if name == "" {
			name = openlane.Deref(n.DisplayName)
		}
		items = append(items, idNameRef{ID: n.ID, Name: name, Status: openlane.Deref(n.EntityRelationshipStateName)})
	}
	return &relSummary[idNameRef]{Count: resp.Entities.TotalCount, Items: items}
}

func (h *handlers) fetchAssets(ctx context.Context, where *graphclient.AssetWhereInput) *relSummary[idNameRef] {
	first := relFirst()
	resp, err := h.api.GetAssets(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]idNameRef, 0, len(resp.Assets.Edges))
	for _, e := range resp.Assets.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, idNameRef{ID: n.ID, Name: n.Name, Status: openlane.Format(n.AssetType)})
	}
	return &relSummary[idNameRef]{Count: resp.Assets.TotalCount, Items: items}
}

func (h *handlers) fetchImplementations(ctx context.Context, where *graphclient.ControlImplementationWhereInput) *relSummary[implementationRef] {
	first := relFirst()
	resp, err := h.api.GetControlImplementations(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]implementationRef, 0, len(resp.ControlImplementations.Edges))
	for _, e := range resp.ControlImplementations.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, implementationRef{ID: n.ID, Status: openlane.Format(n.Status), Verified: n.Verified})
	}
	return &relSummary[implementationRef]{Count: resp.ControlImplementations.TotalCount, Items: items}
}

func (h *handlers) fetchAssessmentResponses(ctx context.Context, assessmentID string) *relSummary[assessmentResponseRef] {
	first := relFirst()
	where := &graphclient.AssessmentResponseWhereInput{AssessmentID: &assessmentID}
	resp, err := h.api.GetAssessmentResponses(ctx, &first, nil, where)
	if err != nil {
		return nil
	}
	items := make([]assessmentResponseRef, 0, len(resp.AssessmentResponses.Edges))
	for _, e := range resp.AssessmentResponses.Edges {
		if e == nil || e.Node == nil {
			continue
		}
		n := e.Node
		items = append(items, assessmentResponseRef{
			ID:     n.ID,
			Status: openlane.Format(n.Status),
			Email:  openlane.Deref(n.Email),
			Due:    openlane.Format(n.DueDate),
		})
	}
	return &relSummary[assessmentResponseRef]{Count: resp.AssessmentResponses.TotalCount, Items: items}
}
