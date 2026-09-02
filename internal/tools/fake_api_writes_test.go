package tools

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/theopenlane/go-client/graphclient"
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
func (f *fakeAPI) CreateRisk(context.Context, graphclient.CreateRiskInput) (*graphclient.CreateRisk, error) {
	return nil, errors.New("unused")
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
