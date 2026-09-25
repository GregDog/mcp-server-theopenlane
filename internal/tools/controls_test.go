package tools

import (
	"context"
	"errors"
	"testing"

	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

type fakeAPI struct {
	controls *graphclient.GetControls
	control  *graphclient.GetControlByID
	tasks    *graphclient.GetTasks
	err      error

	policy     *graphclient.GetInternalPolicyByID
	groups     *graphclient.GetGroups
	group      *graphclient.GetGroupByID
	groupsByID map[string]graphclient.GetGroupByID_Group
	users      *graphclient.GetUsers
	orgMembers *graphclient.GetOrgMembersByOrgID
	user       *graphclient.GetUserByID
	workflow   *graphclient.GetWorkflowDefinitionByID
	createdWF  *graphclient.CreateWorkflowDefinition
	metadata   *openlane.WorkflowMetadata
	assignment *openlane.WorkflowAssignmentDetail
	deletedID  string
	entity     *openlane.EntityDetail

	platforms              *graphclient.GetPlatforms
	platformDetail         *openlane.PlatformDetail
	lastCreatePlatformInput graphclient.CreatePlatformInput
	lastUpdatePlatformInput graphclient.UpdatePlatformInput
	lastPlatformWhere      *graphclient.PlatformWhereInput

	lastCreateEvidenceInput graphclient.CreateEvidenceInput
	lastUpdateEvidenceInput graphclient.UpdateEvidenceInput
	lastUpdateControlInput  graphclient.UpdateControlInput
	lastUpdateRiskInput     graphclient.UpdateRiskInput
	lastUpdatePolicyInput   graphclient.UpdateInternalPolicyInput
	lastControlWhere        *graphclient.ControlWhereInput

	risk  *graphclient.GetRiskByID
	risks *graphclient.GetRisks

	mappedControl  *graphclient.GetMappedControlByID
	mappedControls *graphclient.GetMappedControls

	lastCreateMappedControlInput graphclient.CreateMappedControlInput
	lastUpdateMappedControlInput graphclient.UpdateMappedControlInput
}

func (f *fakeAPI) GetControls(ctx context.Context, first *int64, after *string, where *graphclient.ControlWhereInput) (*graphclient.GetControls, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastControlWhere = where
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
func (f *fakeAPI) GetInternalPolicyByID(_ context.Context, id string) (*graphclient.GetInternalPolicyByID, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.policy != nil {
		return f.policy, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetRisks(_ context.Context, _ *int64, _ *string, _ *graphclient.RiskWhereInput) (*graphclient.GetRisks, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.risks != nil {
		return f.risks, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetRiskByID(_ context.Context, _ string) (*graphclient.GetRiskByID, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.risk != nil {
		return f.risk, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetStandards(context.Context, *int64, *string, *graphclient.StandardWhereInput) (*graphclient.GetStandards, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetStandardByID(context.Context, string) (*graphclient.GetStandardByID, error) {
	return nil, errors.New("unused")
}

func (f *fakeAPI) GetTasks(ctx context.Context, first *int64, after *string, where *graphclient.TaskWhereInput) (*graphclient.GetTasks, error) {
	if f.err != nil {
		return nil, f.err
	}
	return f.tasks, nil
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
	_, page, err := h.listControls(context.Background(), nil, controlListInput{Limit: 10})
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
	_, _, err := h.searchControls(context.Background(), nil, controlSearchInput{})
	if err == nil {
		t.Fatal("expected error")
	}
}

func TestSearchControlsLinkableOnly(t *testing.T) {
	api := &fakeAPI{controls: &graphclient.GetControls{Controls: graphclient.GetControls_Controls{}}}
	h := &handlers{api: api}
	_, _, err := h.searchControls(context.Background(), nil, controlSearchInput{
		Query:        "awareness",
		LinkableOnly: true,
	})
	if err != nil {
		t.Fatal(err)
	}
	if api.lastControlWhere == nil || api.lastControlWhere.OwnerIDNotNil == nil || !*api.lastControlWhere.OwnerIDNotNil {
		t.Fatalf("expected OwnerIDNotNil filter: %+v", api.lastControlWhere)
	}
}
