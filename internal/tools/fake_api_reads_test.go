package tools

import (
	"context"
	"errors"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func (f *fakeAPI) GetTaskByID(context.Context, string) (*graphclient.GetTaskByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetEntities(context.Context, *int64, *string, *graphclient.EntityWhereInput) (*graphclient.GetEntities, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetEntityByID(context.Context, string) (*graphclient.GetEntityByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetEntityDetail(_ context.Context, id string) (*openlane.EntityDetail, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.entity != nil {
		return f.entity, nil
	}
	return &openlane.EntityDetail{ID: id, Name: strPtr("Acme Vendor")}, nil
}
func (f *fakeAPI) GetAssets(context.Context, *int64, *string, *graphclient.AssetWhereInput) (*graphclient.GetAssets, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetAssetByID(context.Context, string) (*graphclient.GetAssetByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetContacts(context.Context, *int64, *string, *graphclient.ContactWhereInput) (*graphclient.GetContacts, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetContactByID(context.Context, string) (*graphclient.GetContactByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetReviews(context.Context, *int64, *string, *graphclient.ReviewWhereInput) (*graphclient.GetReviews, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetReviewByID(context.Context, string) (*graphclient.GetReviewByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetControlImplementations(context.Context, *int64, *string, *graphclient.ControlImplementationWhereInput) (*graphclient.GetControlImplementations, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetControlImplementationByID(context.Context, string) (*graphclient.GetControlImplementationByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetAssessments(context.Context, *int64, *string, *graphclient.AssessmentWhereInput) (*graphclient.GetAssessments, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetAssessmentByID(context.Context, string) (*graphclient.GetAssessmentByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetAssessmentResponses(context.Context, *int64, *string, *graphclient.AssessmentResponseWhereInput) (*graphclient.GetAssessmentResponses, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetFindings(context.Context, *int64, *string, *graphclient.FindingWhereInput) (*graphclient.GetFindings, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetFindingByID(context.Context, string) (*graphclient.GetFindingByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetRemediations(context.Context, *int64, *string, *graphclient.RemediationWhereInput) (*graphclient.GetRemediations, error) {
	return nil, errors.New("unused")
}
