package openlane

import (
	"context"
	"fmt"
	"time"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"
)

// WorkflowMetadata is the workflowMetadata query result (no go-client method in v0.14.0).
type WorkflowMetadata struct {
	ObjectTypes []WorkflowObjectTypeMetadata `json:"objectTypes"`
	Extensions  map[string]any               `json:"extensions"`
}

// WorkflowObjectTypeMetadata describes a workflow-eligible object type.
type WorkflowObjectTypeMetadata struct {
	Type           string                  `json:"type"`
	Label          string                  `json:"label"`
	Description    string                  `json:"description"`
	EligibleFields []WorkflowFieldMetadata `json:"eligibleFields"`
	EligibleEdges  []string                `json:"eligibleEdges"`
	ResolverKeys   []string                `json:"resolverKeys"`
}

// WorkflowFieldMetadata describes a workflow-eligible field.
type WorkflowFieldMetadata struct {
	Name  string `json:"name"`
	Label string `json:"label"`
	Type  string `json:"type"`
}

// WorkflowAssignmentDetail adds fields missing from go-client GetWorkflowAssignmentByID (v0.14.0).
type WorkflowAssignmentDetail struct {
	ID                       string
	DisplayID                string
	AssignmentKey            string
	Status                   string
	Role                     string
	Label                    string
	Required                 bool
	Notes                    string
	DueAt                    string
	WorkflowInstanceID       string
	ActorUserID              string
	ActorGroupID             string
	CreatedBy                string
	DecidedAt                string
	Targets                  []WorkflowAssignmentTarget
	InstanceWorkflowDefID    string
	InstanceState            string
	InstanceInternalPolicyID string
	InstanceControlID        string
	InstanceEvidenceID       string
}

// WorkflowAssignmentTarget is an approval target on an assignment.
type WorkflowAssignmentTarget struct {
	ID            string `json:"id"`
	TargetType    string `json:"target_type"`
	TargetUserID  string `json:"target_user_id,omitempty"`
	TargetGroupID string `json:"target_group_id,omitempty"`
	ResolverKey   string `json:"resolver_key,omitempty"`
}

// WorkflowInstanceDetail adds fields missing from go-client GetWorkflowInstanceByID (v0.14.0).
type WorkflowInstanceDetail struct {
	CurrentActionIndex *int64
	ProposalPreview    *WorkflowProposalPreview
}

// WorkflowProposalPreview is the instance proposal preview when available.
type WorkflowProposalPreview struct {
	ProposalID        string              `json:"proposal_id"`
	DomainKey         string              `json:"domain_key"`
	State             string              `json:"state"`
	SubmittedAt       string              `json:"submitted_at,omitempty"`
	SubmittedByUserID string              `json:"submitted_by_user_id,omitempty"`
	ProposedChanges   map[string]any      `json:"proposed_changes,omitempty"`
	CurrentValues     map[string]any      `json:"current_values,omitempty"`
	Diffs             []WorkflowFieldDiff `json:"diffs,omitempty"`
}

// WorkflowFieldDiff is a single field change in a proposal preview.
type WorkflowFieldDiff struct {
	Field         string `json:"field"`
	Label         string `json:"label,omitempty"`
	Type          string `json:"type,omitempty"`
	CurrentValue  any    `json:"current_value,omitempty"`
	ProposedValue any    `json:"proposed_value,omitempty"`
	Diff          string `json:"diff,omitempty"`
}

const workflowMetadataQuery = `query WorkflowMetadata {
  workflowMetadata {
    objectTypes {
      type
      label
      description
      eligibleFields { name label type }
      eligibleEdges
      resolverKeys
    }
    extensions
  }
}`

const workflowAssignmentDetailQuery = `query WorkflowAssignmentDetail($id: ID!) {
  workflowAssignment(id: $id) {
    id
    displayID
    assignmentKey
    status
    role
    label
    required
    notes
    dueAt
    workflowInstanceID
    actorUserID
    actorGroupID
    createdBy
    decidedAt
    workflowAssignmentTargets {
      edges {
        node {
          id
          targetType
          targetUserID
          targetGroupID
          resolverKey
        }
      }
    }
    workflowInstance {
      workflowDefinitionID
      state
      internalPolicyID
      controlID
      evidenceID
    }
  }
}`

const workflowInstanceDetailQuery = `query WorkflowInstanceDetail($id: ID!) {
  workflowInstance(id: $id) {
    currentActionIndex
    proposalPreview {
      proposalID
      domainKey
      state
      submittedAt
      submittedByUserID
      proposedChanges
      currentValues
      diffs {
        field
        label
        type
        currentValue
        proposedValue
        diff
      }
    }
  }
}`

type workflowMetadataResponse struct {
	WorkflowMetadata WorkflowMetadata `json:"workflowMetadata"`
}

type workflowAssignmentDetailResponse struct {
	WorkflowAssignment struct {
		ID                        string                          `json:"id"`
		DisplayID                 string                          `json:"displayID"`
		AssignmentKey             string                          `json:"assignmentKey"`
		Status                    *enums.WorkflowAssignmentStatus `json:"status"`
		Role                      string                          `json:"role"`
		Label                     *string                         `json:"label"`
		Required                  bool                            `json:"required"`
		Notes                     *string                         `json:"notes"`
		DueAt                     *time.Time                      `json:"dueAt"`
		WorkflowInstanceID        string                          `json:"workflowInstanceID"`
		ActorUserID               *string                         `json:"actorUserID"`
		ActorGroupID              *string                         `json:"actorGroupID"`
		CreatedBy                 *string                         `json:"createdBy"`
		DecidedAt                 *time.Time                      `json:"decidedAt"`
		WorkflowAssignmentTargets struct {
			Edges []struct {
				Node *struct {
					ID            string  `json:"id"`
					TargetType    string  `json:"targetType"`
					TargetUserID  *string `json:"targetUserID"`
					TargetGroupID *string `json:"targetGroupID"`
					ResolverKey   *string `json:"resolverKey"`
				} `json:"node"`
			} `json:"edges"`
		} `json:"workflowAssignmentTargets"`
		WorkflowInstance *struct {
			WorkflowDefinitionID string                      `json:"workflowDefinitionID"`
			State                enums.WorkflowInstanceState `json:"state"`
			InternalPolicyID     *string                     `json:"internalPolicyID"`
			ControlID            *string                     `json:"controlID"`
			EvidenceID           *string                     `json:"evidenceID"`
		} `json:"workflowInstance"`
	} `json:"workflowAssignment"`
}

type workflowInstanceDetailResponse struct {
	WorkflowInstance struct {
		CurrentActionIndex *int64 `json:"currentActionIndex"`
		ProposalPreview    *struct {
			ProposalID        string         `json:"proposalID"`
			DomainKey         string         `json:"domainKey"`
			State             string         `json:"state"`
			SubmittedAt       *time.Time     `json:"submittedAt"`
			SubmittedByUserID *string        `json:"submittedByUserID"`
			ProposedChanges   map[string]any `json:"proposedChanges"`
			CurrentValues     map[string]any `json:"currentValues"`
			Diffs             []struct {
				Field         string  `json:"field"`
				Label         *string `json:"label"`
				Type          *string `json:"type"`
				CurrentValue  any     `json:"currentValue"`
				ProposedValue any     `json:"proposedValue"`
				Diff          *string `json:"diff"`
			} `json:"diffs"`
		} `json:"proposalPreview"`
	} `json:"workflowInstance"`
}

func (a *api) graphClient() (*graphclient.Client, error) {
	gc, ok := a.c.GraphClient.(*graphclient.Client)
	if !ok {
		return nil, fmt.Errorf("graph client type assertion failed")
	}
	return gc, nil
}

func (a *api) GetWorkflowMetadata(ctx context.Context) (*WorkflowMetadata, error) {
	gc, err := a.graphClient()
	if err != nil {
		return nil, err
	}
	var res workflowMetadataResponse
	if err := gc.Client.Post(ctx, "WorkflowMetadata", workflowMetadataQuery, &res, nil); err != nil {
		return nil, RedactError(err)
	}
	return &res.WorkflowMetadata, nil
}

func (a *api) GetWorkflowAssignmentDetail(ctx context.Context, id string) (*WorkflowAssignmentDetail, error) {
	gc, err := a.graphClient()
	if err != nil {
		return nil, err
	}
	var res workflowAssignmentDetailResponse
	vars := map[string]any{"id": id}
	if err := gc.Client.Post(ctx, "WorkflowAssignmentDetail", workflowAssignmentDetailQuery, &res, vars); err != nil {
		return nil, RedactError(err)
	}
	wa := res.WorkflowAssignment
	out := &WorkflowAssignmentDetail{
		ID:                 wa.ID,
		DisplayID:          wa.DisplayID,
		AssignmentKey:      wa.AssignmentKey,
		Status:             Format(wa.Status),
		Role:               wa.Role,
		Label:              Deref(wa.Label),
		Required:           wa.Required,
		Notes:              Deref(wa.Notes),
		DueAt:              Format(wa.DueAt),
		WorkflowInstanceID: wa.WorkflowInstanceID,
		ActorUserID:        Deref(wa.ActorUserID),
		ActorGroupID:       Deref(wa.ActorGroupID),
		CreatedBy:          Deref(wa.CreatedBy),
		DecidedAt:          Format(wa.DecidedAt),
	}
	for _, e := range wa.WorkflowAssignmentTargets.Edges {
		if e.Node == nil {
			continue
		}
		n := e.Node
		out.Targets = append(out.Targets, WorkflowAssignmentTarget{
			ID:            n.ID,
			TargetType:    n.TargetType,
			TargetUserID:  Deref(n.TargetUserID),
			TargetGroupID: Deref(n.TargetGroupID),
			ResolverKey:   Deref(n.ResolverKey),
		})
	}
	if wa.WorkflowInstance != nil {
		inst := wa.WorkflowInstance
		out.InstanceWorkflowDefID = inst.WorkflowDefinitionID
		out.InstanceState = inst.State.String()
		out.InstanceInternalPolicyID = Deref(inst.InternalPolicyID)
		out.InstanceControlID = Deref(inst.ControlID)
		out.InstanceEvidenceID = Deref(inst.EvidenceID)
	}
	return out, nil
}

func (a *api) GetWorkflowInstanceDetail(ctx context.Context, id string) (*WorkflowInstanceDetail, error) {
	gc, err := a.graphClient()
	if err != nil {
		return nil, err
	}
	var res workflowInstanceDetailResponse
	vars := map[string]any{"id": id}
	if err := gc.Client.Post(ctx, "WorkflowInstanceDetail", workflowInstanceDetailQuery, &res, vars); err != nil {
		return nil, RedactError(err)
	}
	wi := res.WorkflowInstance
	out := &WorkflowInstanceDetail{
		CurrentActionIndex: wi.CurrentActionIndex,
	}
	if wi.ProposalPreview != nil {
		pp := wi.ProposalPreview
		preview := &WorkflowProposalPreview{
			ProposalID:        pp.ProposalID,
			DomainKey:         pp.DomainKey,
			State:             pp.State,
			SubmittedAt:       Format(pp.SubmittedAt),
			SubmittedByUserID: Deref(pp.SubmittedByUserID),
			ProposedChanges:   pp.ProposedChanges,
			CurrentValues:     pp.CurrentValues,
		}
		for _, d := range pp.Diffs {
			preview.Diffs = append(preview.Diffs, WorkflowFieldDiff{
				Field:         d.Field,
				Label:         Deref(d.Label),
				Type:          Deref(d.Type),
				CurrentValue:  d.CurrentValue,
				ProposedValue: d.ProposedValue,
				Diff:          Deref(d.Diff),
			})
		}
		out.ProposalPreview = preview
	}
	return out, nil
}

// requestChangesWorkflowAssignment and reassignWorkflowAssignment exist in the
// live GraphQL schema (core workflowassignmentextended.graphql) but have no
// go-client methods in v0.14.0.
const requestChangesWorkflowAssignmentMutation = `mutation RequestChangesWorkflowAssignment($id: ID!, $reason: String, $inputs: Map) {
  requestChangesWorkflowAssignment(id: $id, reason: $reason, inputs: $inputs) {
    workflowAssignment { id displayID status workflowInstanceID }
  }
}`

const reassignWorkflowAssignmentMutation = `mutation ReassignWorkflowAssignment($id: ID!, $targetUserID: ID!) {
  reassignWorkflowAssignment(id: $id, targetUserID: $targetUserID) {
    id
    displayID
    status
    workflowInstanceID
  }
}`

type requestChangesWorkflowAssignmentResponse struct {
	RequestChangesWorkflowAssignment struct {
		WorkflowAssignment struct {
			ID string `json:"id"`
		} `json:"workflowAssignment"`
	} `json:"requestChangesWorkflowAssignment"`
}

type reassignWorkflowAssignmentResponse struct {
	ReassignWorkflowAssignment struct {
		ID string `json:"id"`
	} `json:"reassignWorkflowAssignment"`
}

func (a *api) RequestChangesWorkflowAssignment(ctx context.Context, id string, reason *string, inputs map[string]any) error {
	gc, err := a.graphClient()
	if err != nil {
		return err
	}
	vars := map[string]any{"id": id, "reason": reason, "inputs": inputs}
	var res requestChangesWorkflowAssignmentResponse
	if err := gc.Client.Post(ctx, "RequestChangesWorkflowAssignment", requestChangesWorkflowAssignmentMutation, &res, vars); err != nil {
		return RedactError(err)
	}
	return nil
}

func (a *api) ReassignWorkflowAssignment(ctx context.Context, id, targetUserID string) (string, error) {
	gc, err := a.graphClient()
	if err != nil {
		return "", err
	}
	vars := map[string]any{"id": id, "targetUserID": targetUserID}
	var res reassignWorkflowAssignmentResponse
	if err := gc.Client.Post(ctx, "ReassignWorkflowAssignment", reassignWorkflowAssignmentMutation, &res, vars); err != nil {
		return "", RedactError(err)
	}
	return res.ReassignWorkflowAssignment.ID, nil
}
