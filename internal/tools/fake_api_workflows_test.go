package tools

import (
	"context"
	"errors"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func (f *fakeAPI) GetWorkflowDefinitions(context.Context, *int64, *string, *graphclient.WorkflowDefinitionWhereInput) (*graphclient.GetWorkflowDefinitions, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowInstances(context.Context, *int64, *string, *graphclient.WorkflowInstanceWhereInput) (*graphclient.GetWorkflowInstances, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowInstanceByID(context.Context, string) (*graphclient.GetWorkflowInstanceByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowEvents(context.Context, *int64, *string, *graphclient.WorkflowEventWhereInput) (*graphclient.GetWorkflowEvents, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetMyWorkflowAssignments(context.Context, *int64, *string, *graphclient.WorkflowAssignmentWhereInput) (*graphclient.GetMyWorkflowAssignments, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowAssignments(context.Context, *int64, *string, *graphclient.WorkflowAssignmentWhereInput) (*graphclient.GetWorkflowAssignments, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowAssignmentByID(context.Context, string) (*graphclient.GetWorkflowAssignmentByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowInstanceDetail(context.Context, string) (*openlane.WorkflowInstanceDetail, error) {
	return nil, errors.New("unused")
}

func (f *fakeAPI) GetGroups(_ context.Context, _ *int64, _ *string, _ *graphclient.GroupWhereInput) (*graphclient.GetGroups, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.groups != nil {
		return f.groups, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetGroupByID(_ context.Context, _ string) (*graphclient.GetGroupByID, error) {
	if f.group != nil {
		return f.group, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetUsers(_ context.Context, _ *int64, _ *string, _ *graphclient.UserWhereInput) (*graphclient.GetUsers, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.users != nil {
		return f.users, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetUserByID(_ context.Context, _ string) (*graphclient.GetUserByID, error) {
	if f.user != nil {
		return f.user, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowMetadata(context.Context) (*openlane.WorkflowMetadata, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.metadata != nil {
		return f.metadata, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetWorkflowDefinitionByID(_ context.Context, id string) (*graphclient.GetWorkflowDefinitionByID, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.deletedID != "" && id == f.deletedID {
		return nil, errors.New("not found")
	}
	if f.workflow != nil {
		return f.workflow, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) CreateWorkflowDefinition(_ context.Context, input graphclient.CreateWorkflowDefinitionInput) (*graphclient.CreateWorkflowDefinition, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.createdWF != nil {
		return f.createdWF, nil
	}
	id := "wfd-created"
	if f.workflow == nil {
		f.workflow = &graphclient.GetWorkflowDefinitionByID{
			WorkflowDefinition: graphclient.GetWorkflowDefinitionByID_WorkflowDefinition{
				ID:             id,
				Name:           input.Name,
				SchemaType:     input.SchemaType,
				WorkflowKind:   input.WorkflowKind,
				Active:         openlane.Deref(input.Active),
				Draft:          openlane.Deref(input.Draft),
				DefinitionJSON: input.DefinitionJSON,
				TrackedFields:  input.TrackedFields,
			},
		}
	}
	return &graphclient.CreateWorkflowDefinition{
		CreateWorkflowDefinition: graphclient.CreateWorkflowDefinition_CreateWorkflowDefinition{
			WorkflowDefinition: graphclient.CreateWorkflowDefinition_CreateWorkflowDefinition_WorkflowDefinition{
				ID:   id,
				Name: input.Name,
			},
		},
	}, nil
}
func (f *fakeAPI) UpdateWorkflowDefinition(_ context.Context, id string, input graphclient.UpdateWorkflowDefinitionInput) (*graphclient.UpdateWorkflowDefinition, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.workflow != nil {
		if input.Name != nil {
			f.workflow.WorkflowDefinition.Name = *input.Name
		}
		if input.Active != nil {
			f.workflow.WorkflowDefinition.Active = *input.Active
		}
		if input.DefinitionJSON != nil {
			f.workflow.WorkflowDefinition.DefinitionJSON = input.DefinitionJSON
		}
	}
	return &graphclient.UpdateWorkflowDefinition{}, nil
}
func (f *fakeAPI) DeleteWorkflowDefinition(_ context.Context, id string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.deletedID = id
	return id, nil
}
func (f *fakeAPI) ApproveWorkflowAssignment(_ context.Context, id string) (*graphclient.ApproveWorkflowAssignment, error) {
	if f.assignment != nil {
		f.assignment.Status = "APPROVED"
	}
	return &graphclient.ApproveWorkflowAssignment{}, nil
}
func (f *fakeAPI) RejectWorkflowAssignment(_ context.Context, id string, _ *string) (*graphclient.RejectWorkflowAssignment, error) {
	if f.assignment != nil {
		f.assignment.Status = "REJECTED"
	}
	return &graphclient.RejectWorkflowAssignment{}, nil
}
func (f *fakeAPI) RequestChangesWorkflowAssignment(_ context.Context, _ string, _ *string, _ map[string]any) error {
	if f.assignment != nil {
		f.assignment.Status = "CHANGES_REQUESTED"
	}
	return nil
}
func (f *fakeAPI) ReassignWorkflowAssignment(_ context.Context, _ string, _ string) (string, error) {
	return "asg-reassigned", nil
}
func (f *fakeAPI) GetWorkflowAssignmentDetail(_ context.Context, id string) (*openlane.WorkflowAssignmentDetail, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.assignment != nil {
		if f.assignment.ID == "" {
			f.assignment.ID = id
		}
		return f.assignment, nil
	}
	return nil, errors.New("unused")
}
