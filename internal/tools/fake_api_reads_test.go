package tools

import (
	"context"
	"errors"

	"github.com/theopenlane/go-client/graphclient"
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
