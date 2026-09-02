package tools

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
)

func summarizeWorkflowDefinition(name, schemaType, kind string, active, draft bool, revision int64, cooldown int64, tracked []string, doc *models.WorkflowDefinitionDocument) string {
	var b strings.Builder
	b.WriteString(fmt.Sprintf("%q is a %s workflow for %s", name, kind, schemaType))
	if active {
		b.WriteString(" (active")
	} else {
		b.WriteString(" (inactive")
	}
	if draft {
		b.WriteString(", draft")
	}
	b.WriteString(")")
	if revision > 0 {
		b.WriteString(fmt.Sprintf("; revision %d", revision))
	}
	if cooldown > 0 {
		b.WriteString(fmt.Sprintf("; cooldown %ds", cooldown))
	}
	b.WriteString(".")

	if doc == nil {
		if len(tracked) > 0 {
			b.WriteString(fmt.Sprintf(" Tracked fields: %s.", strings.Join(tracked, ", ")))
		}
		return b.String()
	}

	if doc.ApprovalTiming != "" {
		b.WriteString(fmt.Sprintf(" Approval timing: %s.", doc.ApprovalTiming))
	}
	if doc.ApprovalSubmissionMode != "" {
		b.WriteString(fmt.Sprintf(" Submission mode: %s.", doc.ApprovalSubmissionMode))
	}

	if len(doc.Triggers) > 0 {
		b.WriteString(" Triggers:")
		for i, t := range doc.Triggers {
			if i >= 5 {
				b.WriteString(fmt.Sprintf(" …and %d more.", len(doc.Triggers)-5))
				break
			}
			parts := []string{t.Operation}
			if len(t.Fields) > 0 {
				parts = append(parts, "fields="+strings.Join(t.Fields, ","))
			}
			if len(t.Edges) > 0 {
				parts = append(parts, "edges="+strings.Join(t.Edges, ","))
			}
			if t.Expression != "" {
				parts = append(parts, "expr="+truncate(t.Expression, 60))
			}
			b.WriteString(" [" + strings.Join(parts, " ") + "]")
		}
		b.WriteString(".")
	}

	if len(doc.Conditions) > 0 {
		b.WriteString(" Conditions:")
		for i, c := range doc.Conditions {
			if i >= 3 {
				break
			}
			if c.Expression != "" {
				b.WriteString(" " + truncate(c.Expression, 80))
			}
		}
		b.WriteString(".")
	}

	if len(doc.Actions) > 0 {
		b.WriteString(" Actions:")
		for i, act := range doc.Actions {
			if i >= 5 {
				b.WriteString(fmt.Sprintf(" …and %d more.", len(doc.Actions)-5))
				break
			}
			b.WriteString(fmt.Sprintf(" %s(%s)", act.Key, act.Type))
			if act.Type == string(enums.WorkflowActionTypeApproval) || act.Type == "REQUEST_APPROVAL" {
				if targets := summarizeApprovalTargets(act.Params); targets != "" {
					b.WriteString(" → " + targets)
				}
			}
			if act.When != "" {
				b.WriteString(" when " + truncate(act.When, 40))
			}
		}
		b.WriteString(".")
	}

	return b.String()
}

func summarizeApprovalTargets(params json.RawMessage) string {
	if len(params) == 0 {
		return ""
	}
	var p struct {
		Targets []struct {
			Type        string `json:"type"`
			ID          string `json:"id"`
			ResolverKey string `json:"resolver_key"`
		} `json:"targets"`
		RequiredCount *int `json:"required_count"`
	}
	if err := json.Unmarshal(params, &p); err != nil {
		return ""
	}
	var parts []string
	for _, t := range p.Targets {
		switch t.Type {
		case "RESOLVER":
			parts = append(parts, "RESOLVER:"+t.ResolverKey)
		case "GROUP":
			parts = append(parts, "GROUP:"+truncateID(t.ID))
		case "USER":
			parts = append(parts, "USER:"+truncateID(t.ID))
		case "ROLE":
			parts = append(parts, "ROLE:"+truncateID(t.ID))
		case "CHANNEL":
			parts = append(parts, "CHANNEL")
		default:
			parts = append(parts, t.Type)
		}
	}
	if p.RequiredCount != nil && *p.RequiredCount > 0 {
		parts = append(parts, fmt.Sprintf("quorum=%d", *p.RequiredCount))
	}
	return strings.Join(parts, ", ")
}

func truncate(s string, max int) string {
	s = strings.TrimSpace(s)
	if len(s) <= max {
		return s
	}
	return s[:max-3] + "..."
}

func truncateID(id string) string {
	if len(id) <= 12 {
		return id
	}
	return id[:8] + "…"
}
