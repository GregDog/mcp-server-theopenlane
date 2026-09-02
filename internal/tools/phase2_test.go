package tools

import (
	"context"
	"encoding/json"
	"strings"
	"testing"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func testWorkflowMetadata() *openlane.WorkflowMetadata {
	return &openlane.WorkflowMetadata{
		ObjectTypes: []openlane.WorkflowObjectTypeMetadata{{
			Type:  "InternalPolicy",
			Label: "Internal Policy",
			EligibleFields: []openlane.WorkflowFieldMetadata{
				{Name: "status", Label: "Status", Type: "enum"},
			},
			EligibleEdges: []string{},
			ResolverKeys:  []string{"POLICY_APPROVER", "POLICY_DELEGATE", "POLICY_OWNER"},
		}},
	}
}

func TestValidatePolicyTransition(t *testing.T) {
	cases := []struct {
		current, target string
		required        bool
		ok              bool
	}{
		{"DRAFT", "NEEDS_APPROVAL", true, true},
		{"NEEDS_APPROVAL", "NEEDS_APPROVAL", true, false},
		{"ARCHIVED", "NEEDS_APPROVAL", true, false},
		{"NEEDS_APPROVAL", "APPROVED", true, true},
		{"DRAFT", "APPROVED", true, false},
		{"APPROVED", "PUBLISHED", true, true},
		{"DRAFT", "PUBLISHED", false, true},
		{"DRAFT", "PUBLISHED", true, false},
		{"NEEDS_APPROVAL", "DRAFT", true, true},
		{"APPROVED", "DRAFT", true, true},
		{"PUBLISHED", "DRAFT", true, false},
	}
	for _, tc := range cases {
		err := validatePolicyTransition(tc.current, tc.target, tc.required)
		if tc.ok && err != nil {
			t.Errorf("%s -> %s required=%v: unexpected error %v", tc.current, tc.target, tc.required, err)
		}
		if !tc.ok && err == nil {
			t.Errorf("%s -> %s required=%v: expected error", tc.current, tc.target, tc.required)
		}
	}
}

func TestPolicyApproveRequiresConfirm(t *testing.T) {
	status := enums.DocumentNeedsApproval
	h := &handlers{api: &fakeAPI{policy: &graphclient.GetInternalPolicyByID{
		InternalPolicy: graphclient.GetInternalPolicyByID_InternalPolicy{
			ID:               "pol-1",
			DisplayID:        "PLC-1",
			Name:             "Clean Desk",
			Status:           &status,
			ApprovalRequired: boolPtr(true),
		},
	}}, allowWrite: true}

	_, out, err := h.approvePolicy(context.Background(), nil, policyLifecycleInput{ID: "pol-1"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Error != errConfirmationRequired {
		t.Fatalf("expected confirmation required, got %+v", out)
	}
	if out.Confirmed {
		t.Fatal("should not confirm")
	}
	if out.CurrentStatus != "NEEDS_APPROVAL" || out.RequestedStatus != "APPROVED" {
		t.Fatalf("preview: %+v", out)
	}
}

func TestPolicyApproveUpdatesStatus(t *testing.T) {
	status := enums.DocumentNeedsApproval
	h := &handlers{api: &fakeAPI{policy: &graphclient.GetInternalPolicyByID{
		InternalPolicy: graphclient.GetInternalPolicyByID_InternalPolicy{
			ID:               "pol-1",
			Name:             "Clean Desk",
			Status:           &status,
			ApprovalRequired: boolPtr(true),
		},
	}}, allowWrite: true}

	_, out, err := h.approvePolicy(context.Background(), nil, policyLifecycleInput{ID: "pol-1", Confirm: true})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Confirmed || out.ResultStatus != "APPROVED" {
		t.Fatalf("got %+v", out)
	}
}

func TestPolicyApproveRefusesWrongStatus(t *testing.T) {
	status := enums.DocumentDraft
	h := &handlers{api: &fakeAPI{policy: &graphclient.GetInternalPolicyByID{
		InternalPolicy: graphclient.GetInternalPolicyByID_InternalPolicy{
			ID:     "pol-1",
			Name:   "Clean Desk",
			Status: &status,
		},
	}}, allowWrite: true}

	_, _, err := h.approvePolicy(context.Background(), nil, policyLifecycleInput{ID: "pol-1", Confirm: true})
	if err == nil || !strings.Contains(err.Error(), "NEEDS_APPROVAL") {
		t.Fatalf("got %v", err)
	}
}

func TestWorkflowAssignmentApproveRequiresPending(t *testing.T) {
	h := &handlers{api: &fakeAPI{assignment: &openlane.WorkflowAssignmentDetail{
		ID:     "asg-1",
		Status: "APPROVED",
	}}, allowWrite: true}
	_, _, err := h.approveWorkflowAssignment(context.Background(), nil, workflowAssignmentActionInput{ID: "asg-1", Confirm: true})
	if err == nil || !strings.Contains(err.Error(), "PENDING") {
		t.Fatalf("got %v", err)
	}
}

func TestWorkflowAssignmentApproveConfirms(t *testing.T) {
	h := &handlers{api: &fakeAPI{assignment: &openlane.WorkflowAssignmentDetail{
		ID:                 "asg-1",
		DisplayID:          "WAS-1",
		Status:             "PENDING",
		WorkflowInstanceID: "inst-1",
	}}, allowWrite: true}

	_, out, err := h.approveWorkflowAssignment(context.Background(), nil, workflowAssignmentActionInput{ID: "asg-1"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Error != errConfirmationRequired {
		t.Fatalf("preview: %+v", out)
	}

	_, out, err = h.approveWorkflowAssignment(context.Background(), nil, workflowAssignmentActionInput{ID: "asg-1", Confirm: true})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Confirmed || out.ResultStatus != "APPROVED" {
		t.Fatalf("got %+v", out)
	}
}

func TestRejectAssignmentRequiresReason(t *testing.T) {
	h := &handlers{api: &fakeAPI{assignment: &openlane.WorkflowAssignmentDetail{ID: "asg-1", Status: "PENDING"}}, allowWrite: true}
	_, _, err := h.rejectWorkflowAssignment(context.Background(), nil, workflowAssignmentActionInput{ID: "asg-1", Confirm: true})
	if err == nil || !strings.Contains(err.Error(), "reason") {
		t.Fatalf("got %v", err)
	}
}

func TestCreateWorkflowRequiresConfirmAndValidates(t *testing.T) {
	h := &handlers{api: &fakeAPI{metadata: testWorkflowMetadata()}, allowWrite: true}
	in := createWorkflowInput{
		Name:         "Policy publication approval",
		SchemaType:   "InternalPolicy",
		Triggers:     []workflowTriggerSpec{{Operation: "UPDATE", Fields: []string{"status"}}},
		ConditionCEL: `object.status == "PUBLISHED"`,
		Actions: []workflowActionSpec{{
			Type:    "REQUEST_APPROVAL",
			Targets: []workflowTargetSpec{{Type: "RESOLVER", ResolverKey: "POLICY_APPROVER"}},
		}},
	}
	_, out, err := h.createWorkflow(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if out.Error != errConfirmationRequired {
		t.Fatalf("preview: %+v", out)
	}

	in.Confirm = true
	_, out, err = h.createWorkflow(context.Background(), nil, in)
	if err != nil {
		t.Fatal(err)
	}
	if !out.Confirmed || out.After == nil || out.After.Name != in.Name {
		t.Fatalf("created: %+v", out)
	}
}

func TestCreateWorkflowRejectsUnknownResolver(t *testing.T) {
	h := &handlers{api: &fakeAPI{metadata: testWorkflowMetadata()}, allowWrite: true}
	_, _, err := h.createWorkflow(context.Background(), nil, createWorkflowInput{
		Name:       "x",
		SchemaType: "InternalPolicy",
		Triggers:   []workflowTriggerSpec{{Operation: "UPDATE", Fields: []string{"status"}}},
		Actions: []workflowActionSpec{{
			Type:    "REQUEST_APPROVAL",
			Targets: []workflowTargetSpec{{Type: "RESOLVER", ResolverKey: "NOT_A_RESOLVER"}},
		}},
		Confirm: true,
	})
	if err == nil || !strings.Contains(err.Error(), "resolver_key") {
		t.Fatalf("got %v", err)
	}
}

func TestCreateWorkflowRejectsUnknownSchemaType(t *testing.T) {
	h := &handlers{api: &fakeAPI{metadata: testWorkflowMetadata()}, allowWrite: true}
	_, _, err := h.createWorkflow(context.Background(), nil, createWorkflowInput{
		Name:       "x",
		SchemaType: "UnknownType",
		Triggers:   []workflowTriggerSpec{{Operation: "UPDATE", Fields: []string{"status"}}},
		Actions:    []workflowActionSpec{{Type: "NOTIFY"}},
		Confirm:    true,
	})
	if err == nil || !strings.Contains(err.Error(), "schema_type") {
		t.Fatalf("got %v", err)
	}
}

func TestUpdateWorkflowActivate(t *testing.T) {
	doc := &models.WorkflowDefinitionDocument{
		Name:         "Policy publication approval",
		SchemaType:   "InternalPolicy",
		WorkflowKind: enums.WorkflowKindApproval,
		Triggers:     []models.WorkflowTrigger{{Operation: "UPDATE", Fields: []string{"status"}}},
		Actions:      []models.WorkflowAction{{Key: "approve", Type: "REQUEST_APPROVAL", Params: json.RawMessage(`{"targets":[{"type":"RESOLVER","resolver_key":"POLICY_APPROVER"}]}`)}},
	}
	h := &handlers{api: &fakeAPI{
		metadata: testWorkflowMetadata(),
		workflow: &graphclient.GetWorkflowDefinitionByID{
			WorkflowDefinition: graphclient.GetWorkflowDefinitionByID_WorkflowDefinition{
				ID:             "wfd-1",
				Name:           doc.Name,
				SchemaType:     "InternalPolicy",
				WorkflowKind:   enums.WorkflowKindApproval,
				Active:         false,
				DefinitionJSON: doc,
			},
		},
	}, allowWrite: true}

	active := true
	_, out, err := h.updateWorkflow(context.Background(), nil, updateWorkflowInput{ID: "wfd-1", Active: &active, Confirm: true})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Confirmed || out.After == nil || !out.After.Active {
		t.Fatalf("got %+v", out)
	}
}

func TestDeleteWorkflowRequiresConfirm(t *testing.T) {
	h := &handlers{api: &fakeAPI{workflow: &graphclient.GetWorkflowDefinitionByID{
		WorkflowDefinition: graphclient.GetWorkflowDefinitionByID_WorkflowDefinition{
			ID:   "wfd-1",
			Name: "To delete",
		},
	}}}
	_, out, err := h.deleteWorkflow(context.Background(), nil, deleteWorkflowInput{ID: "wfd-1"})
	if err != nil {
		t.Fatal(err)
	}
	if out.Error != errConfirmationRequired {
		t.Fatalf("got %+v", out)
	}
	_, out, err = h.deleteWorkflow(context.Background(), nil, deleteWorkflowInput{ID: "wfd-1", Confirm: true})
	if err != nil {
		t.Fatal(err)
	}
	if !out.Confirmed || out.DeletedID != "wfd-1" {
		t.Fatalf("got %+v", out)
	}
}

func TestBuildGroupWhere(t *testing.T) {
	where := buildGroupWhere("Privacy")
	if where == nil || len(where.Or) != 2 {
		t.Fatalf("got %+v", where)
	}
}

func TestBuildUserWhere(t *testing.T) {
	where := buildUserWhere(userListInput{Email: "a@example.com"})
	if where == nil || where.EmailContainsFold == nil {
		t.Fatalf("got %+v", where)
	}
}

func TestEmbeddedWorkflowSchemaPresent(t *testing.T) {
	if len(openlane.WorkflowDefinitionSchemaJSON) < 100 {
		t.Fatal("expected embedded workflow.definition.json")
	}
	if !strings.Contains(string(openlane.WorkflowDefinitionSchemaJSON), "REQUEST_APPROVAL") {
		t.Fatal("schema missing REQUEST_APPROVAL")
	}
}

func boolPtr(v bool) *bool { return &v }
