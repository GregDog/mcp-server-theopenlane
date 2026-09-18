package tools

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func (f *fakeAPI) CreateControl(context.Context, graphclient.CreateControlInput) (*graphclient.CreateControl, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateControl(context.Context, string, graphclient.UpdateControlInput) (*graphclient.UpdateControl, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) CreateEvidence(context.Context, graphclient.CreateEvidenceInput, []*graphql.Upload) (*graphclient.CreateEvidence, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateEvidence(context.Context, string, graphclient.UpdateEvidenceInput, []*graphql.Upload) (*graphclient.UpdateEvidence, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) CreateInternalPolicy(context.Context, graphclient.CreateInternalPolicyInput) (*graphclient.CreateInternalPolicy, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateInternalPolicy(_ context.Context, _ string, input graphclient.UpdateInternalPolicyInput) (*graphclient.UpdateInternalPolicy, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.policy != nil && input.Status != nil {
		f.policy.InternalPolicy.Status = input.Status
	}
	return &graphclient.UpdateInternalPolicy{}, nil
}
func (f *fakeAPI) CreateRisk(_ context.Context, input graphclient.CreateRiskInput) (*graphclient.CreateRisk, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &graphclient.CreateRisk{
		CreateRisk: graphclient.CreateRisk_CreateRisk{
			Risk: graphclient.CreateRisk_CreateRisk_Risk{
				ID:   "risk_1",
				Name: input.Name,
			},
		},
	}, nil
}
func (f *fakeAPI) UpdateRisk(context.Context, string, graphclient.UpdateRiskInput) (*graphclient.UpdateRisk, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) CreateTask(context.Context, graphclient.CreateTaskInput) (*graphclient.CreateTask, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateTask(context.Context, string, graphclient.UpdateTaskInput) (*graphclient.UpdateTask, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) DeleteControl(context.Context, string) (string, error) {
	return "", errors.New("unused")
}
func (f *fakeAPI) DeleteEvidence(context.Context, string) (string, error) {
	return "", errors.New("unused")
}
func (f *fakeAPI) DeleteInternalPolicy(context.Context, string) (string, error) {
	return "", errors.New("unused")
}
func (f *fakeAPI) DeleteRisk(context.Context, string) (string, error) {
	return "", errors.New("unused")
}
func (f *fakeAPI) DeleteTask(context.Context, string) (string, error) {
	return "", errors.New("unused")
}
func (f *fakeAPI) CreateContact(_ context.Context, input graphclient.CreateContactInput) (*graphclient.CreateContact, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &graphclient.CreateContact{
		CreateContact: graphclient.CreateContact_CreateContact{
			Contact: graphclient.CreateContact_CreateContact_Contact{
				ID:       "contact_1",
				FullName: input.FullName,
				Email:    input.Email,
				Tags:     input.Tags,
			},
		},
	}, nil
}
func (f *fakeAPI) CreateEntity(_ context.Context, input graphclient.CreateEntityInput, _ *string, _ *graphql.Upload) (*openlane.EntityDetail, error) {
	if f.err != nil {
		return nil, f.err
	}
	name := "created"
	if input.Name != nil {
		name = *input.Name
	}
	return &openlane.EntityDetail{ID: "ent_1", Name: &name}, nil
}
func (f *fakeAPI) UpdateEntity(_ context.Context, id string, _ graphclient.UpdateEntityInput, _ *graphql.Upload) (*openlane.EntityDetail, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.entity != nil {
		return f.entity, nil
	}
	return &openlane.EntityDetail{ID: id, Name: strPtr("updated")}, nil
}
