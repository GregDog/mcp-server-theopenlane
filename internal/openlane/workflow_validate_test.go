package openlane

import (
	"testing"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
)

func TestValidateWorkflowDefinition(t *testing.T) {
	meta := &WorkflowMetadata{
		ObjectTypes: []WorkflowObjectTypeMetadata{{
			Type:           "InternalPolicy",
			EligibleFields: []WorkflowFieldMetadata{{Name: "status"}},
			ResolverKeys:   []string{"POLICY_APPROVER"},
		}},
	}
	doc := &models.WorkflowDefinitionDocument{
		Name:         "Policy publication approval",
		SchemaType:   "InternalPolicy",
		WorkflowKind: enums.WorkflowKindApproval,
		Triggers:     []models.WorkflowTrigger{{Operation: "UPDATE", Fields: []string{"status"}}},
		Actions: []models.WorkflowAction{{
			Key:    "approve",
			Type:   string(enums.WorkflowActionTypeApproval),
			Params: []byte(`{"targets":[{"type":"RESOLVER","resolver_key":"POLICY_APPROVER"}]}`),
		}},
	}
	if err := ValidateWorkflowDefinition(doc, meta); err != nil {
		t.Fatal(err)
	}

	bad := *doc
	bad.Actions = []models.WorkflowAction{{Key: "approve", Type: "REQUEST_APPROVAL", Params: []byte(`{"targets":[{"type":"RESOLVER","resolver_key":"NOPE"}]}`)}}
	if err := ValidateWorkflowDefinition(&bad, meta); err == nil {
		t.Fatal("expected unknown resolver error")
	}
}
