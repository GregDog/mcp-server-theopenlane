package tools

import (
	"testing"

	"github.com/theopenlane/core/common/enums"
)

func TestBuildWorkflowDefinitionWhere(t *testing.T) {
	active := true
	draft := false
	where := buildWorkflowDefinitionWhere(workflowListInput{
		Name:         "policy",
		SchemaType:   "InternalPolicy",
		WorkflowKind: "APPROVAL",
		Active:       &active,
		Draft:        &draft,
		TrackedField: "status",
	})
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.NameContainsFold == nil || *where.NameContainsFold != "policy" {
		t.Fatalf("name filter: %+v", where.NameContainsFold)
	}
	if where.SchemaTypeContainsFold == nil || *where.SchemaTypeContainsFold != "InternalPolicy" {
		t.Fatalf("schema filter: %+v", where.SchemaTypeContainsFold)
	}
	if where.WorkflowKind == nil || *where.WorkflowKind != enums.WorkflowKindApproval {
		t.Fatalf("kind filter: %+v", where.WorkflowKind)
	}
	if where.Active == nil || !*where.Active {
		t.Fatalf("active filter: %+v", where.Active)
	}
	if where.TrackedFieldsHas == nil || *where.TrackedFieldsHas != "status" {
		t.Fatalf("tracked field filter: %+v", where.TrackedFieldsHas)
	}
}

func TestBuildWorkflowDefinitionSearchWhere(t *testing.T) {
	where := buildWorkflowDefinitionSearchWhere("approval")
	if where == nil || len(where.Or) != 2 {
		t.Fatalf("expected OR search: %+v", where)
	}
}

func TestBuildWorkflowInstanceWhere(t *testing.T) {
	where, err := buildWorkflowInstanceWhere(workflowInstanceListInput{
		DefinitionID: "wfd-1",
		State:        "PAUSED",
		ObjectType:   "Policy",
		ObjectID:     "pol-1",
	})
	if err != nil {
		t.Fatal(err)
	}
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.WorkflowDefinitionID == nil || *where.WorkflowDefinitionID != "wfd-1" {
		t.Fatalf("definition filter: %+v", where.WorkflowDefinitionID)
	}
	if where.State == nil || *where.State != enums.WorkflowInstanceStatePaused {
		t.Fatalf("state filter: %+v", where.State)
	}
	if where.InternalPolicyID == nil || *where.InternalPolicyID != "pol-1" {
		t.Fatalf("policy filter: %+v", where.InternalPolicyID)
	}
}

func TestBuildWorkflowInstanceWhereObjectTypeOnly(t *testing.T) {
	where, err := buildWorkflowInstanceWhere(workflowInstanceListInput{
		ObjectType: "InternalPolicy",
	})
	if err != nil {
		t.Fatal(err)
	}
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.InternalPolicyIDNotNil == nil || !*where.InternalPolicyIDNotNil {
		t.Fatalf("expected internal policy not-nil filter: %+v", where.InternalPolicyIDNotNil)
	}
}

func TestBuildWorkflowInstanceWhereRequiresObjectType(t *testing.T) {
	_, err := buildWorkflowInstanceWhere(workflowInstanceListInput{ObjectID: "x"})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestBuildWorkflowAssignmentWhere(t *testing.T) {
	where := buildWorkflowAssignmentWhere(workflowAssignmentListInput{
		Status:     "PENDING",
		InstanceID: "wfi-1",
	})
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.Status == nil || *where.Status != enums.WorkflowAssignmentStatusPending {
		t.Fatalf("status filter: %+v", where.Status)
	}
	if where.WorkflowInstanceID == nil || *where.WorkflowInstanceID != "wfi-1" {
		t.Fatalf("instance filter: %+v", where.WorkflowInstanceID)
	}
}

func TestSummarizeWorkflowDefinition(t *testing.T) {
	s := summarizeWorkflowDefinition(
		"Policy Approval",
		"InternalPolicy",
		"APPROVAL",
		true, false, 2, 60,
		[]string{"status"},
		nil,
	)
	if s == "" {
		t.Fatal("expected summary")
	}
}
