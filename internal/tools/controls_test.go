package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

type fakeAPI struct {
	controls *graphclient.GetControls
	control  *graphclient.GetControlByID
	err      error
}

func (f *fakeAPI) GetControls(ctx context.Context, first *int64, after *string, where *graphclient.ControlWhereInput) (*graphclient.GetControls, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.controls, nil
}

func (f *fakeAPI) GetControlByID(ctx context.Context, id string) (*graphclient.GetControlByID, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.control, nil
}

func (f *fakeAPI) GetPrograms(context.Context, *int64, *string, *graphclient.ProgramWhereInput) (*graphclient.GetPrograms, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetProgramByID(context.Context, string) (*graphclient.GetProgramByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetEvidences(context.Context, *int64, *string, *graphclient.EvidenceWhereInput) (*graphclient.GetEvidences, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetEvidenceByID(context.Context, string) (*graphclient.GetEvidenceByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetInternalPolicies(context.Context, *int64, *string, *graphclient.InternalPolicyWhereInput) (*graphclient.GetInternalPolicies, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetInternalPolicyByID(context.Context, string) (*graphclient.GetInternalPolicyByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetRisks(context.Context, *int64, *string, *graphclient.RiskWhereInput) (*graphclient.GetRisks, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetRiskByID(context.Context, string) (*graphclient.GetRiskByID, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetStandards(context.Context, *int64, *string, *graphclient.StandardWhereInput) (*graphclient.GetStandards, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetStandardByID(context.Context, string) (*graphclient.GetStandardByID, error) {
	return nil, errors.New("unused")
}

func TestListControls(t *testing.T) {
	title := "Access Control"
	api := &fakeAPI{
		controls: &graphclient.GetControls{
			Controls: graphclient.GetControls_Controls{
				TotalCount: 1,
				PageInfo: graphclient.GetControls_Controls_PageInfo{
					HasNextPage: false,
				},
				Edges: []*graphclient.GetControls_Controls_Edges{
					{
						Node: &graphclient.GetControls_Controls_Edges_Node{
							ID:      "ctrl_1",
							RefCode: "AC-1",
							Title:   &title,
						},
					},
				},
			},
		},
	}
	h := &handlers{api: api}
	_, page, err := h.listControls(context.Background(), nil, listInput{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if page.TotalCount != 1 || len(page.Items) != 1 {
		t.Fatalf("unexpected page: %+v", page)
	}
	if page.Items[0].RefCode != "AC-1" {
		t.Fatalf("ref code: %q", page.Items[0].RefCode)
	}
}

func TestGetControlRequiresID(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.getControl(context.Background(), nil, getInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestGetControl(t *testing.T) {
	title := "Access Control"
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      "ctrl_1",
				RefCode: "AC-1",
				Title:   &title,
			},
		},
	}
	h := &handlers{api: api}
	_, item, err := h.getControl(context.Background(), nil, getInput{ID: "ctrl_1"})
	if err != nil {
		t.Fatal(err)
	}
	if item.RefCode != "AC-1" || item.Title != title {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestSearchControlsRequiresQuery(t *testing.T) {
	h := &handlers{api: &fakeAPI{}}
	_, _, err := h.searchControls(context.Background(), nil, searchInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}
