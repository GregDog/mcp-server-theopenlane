package openlane

import (
	"context"
	"fmt"
	"time"

	"github.com/99designs/gqlgen/graphql"
	gqlclient "github.com/theopenlane/go-client"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
)

const (
	httpTimeout       = 30 * time.Second
	uploadHTTPTimeout = 2 * time.Minute
)

// GraphAPI is the subset of the official Openlane client used by MCP tools.
type GraphAPI interface {
	GetControls(ctx context.Context, first *int64, after *string, where *graphclient.ControlWhereInput) (*graphclient.GetControls, error)
	GetControlByID(ctx context.Context, id string) (*graphclient.GetControlByID, error)
	GetPrograms(ctx context.Context, first *int64, after *string, where *graphclient.ProgramWhereInput) (*graphclient.GetPrograms, error)
	GetProgramByID(ctx context.Context, id string) (*graphclient.GetProgramByID, error)
	GetEvidences(ctx context.Context, first *int64, after *string, where *graphclient.EvidenceWhereInput) (*graphclient.GetEvidences, error)
	GetEvidenceByID(ctx context.Context, id string) (*graphclient.GetEvidenceByID, error)
	GetInternalPolicies(ctx context.Context, first *int64, after *string, where *graphclient.InternalPolicyWhereInput) (*graphclient.GetInternalPolicies, error)
	GetInternalPolicyByID(ctx context.Context, id string) (*graphclient.GetInternalPolicyByID, error)
	GetRisks(ctx context.Context, first *int64, after *string, where *graphclient.RiskWhereInput) (*graphclient.GetRisks, error)
	GetRiskByID(ctx context.Context, id string) (*graphclient.GetRiskByID, error)
	GetStandards(ctx context.Context, first *int64, after *string, where *graphclient.StandardWhereInput) (*graphclient.GetStandards, error)
	GetStandardByID(ctx context.Context, id string) (*graphclient.GetStandardByID, error)
	GetTasks(ctx context.Context, first *int64, after *string, where *graphclient.TaskWhereInput) (*graphclient.GetTasks, error)
	GetTaskByID(ctx context.Context, id string) (*graphclient.GetTaskByID, error)
	GetEntities(ctx context.Context, first *int64, after *string, where *graphclient.EntityWhereInput) (*graphclient.GetEntities, error)
	GetEntityByID(ctx context.Context, id string) (*graphclient.GetEntityByID, error)
	GetAssets(ctx context.Context, first *int64, after *string, where *graphclient.AssetWhereInput) (*graphclient.GetAssets, error)
	GetAssetByID(ctx context.Context, id string) (*graphclient.GetAssetByID, error)
	GetContacts(ctx context.Context, first *int64, after *string, where *graphclient.ContactWhereInput) (*graphclient.GetContacts, error)
	GetContactByID(ctx context.Context, id string) (*graphclient.GetContactByID, error)
	GetControlImplementations(ctx context.Context, first *int64, after *string, where *graphclient.ControlImplementationWhereInput) (*graphclient.GetControlImplementations, error)
	GetControlImplementationByID(ctx context.Context, id string) (*graphclient.GetControlImplementationByID, error)
	GetAssessments(ctx context.Context, first *int64, after *string, where *graphclient.AssessmentWhereInput) (*graphclient.GetAssessments, error)
	GetAssessmentByID(ctx context.Context, id string) (*graphclient.GetAssessmentByID, error)
	GetAssessmentResponses(ctx context.Context, first *int64, after *string, where *graphclient.AssessmentResponseWhereInput) (*graphclient.GetAssessmentResponses, error)
	GetFindings(ctx context.Context, first *int64, after *string, where *graphclient.FindingWhereInput) (*graphclient.GetFindings, error)
	GetFindingByID(ctx context.Context, id string) (*graphclient.GetFindingByID, error)
	GetRemediations(ctx context.Context, first *int64, after *string, where *graphclient.RemediationWhereInput) (*graphclient.GetRemediations, error)
	CreateControl(ctx context.Context, input graphclient.CreateControlInput) (*graphclient.CreateControl, error)
	UpdateControl(ctx context.Context, id string, input graphclient.UpdateControlInput) (*graphclient.UpdateControl, error)
	CreateEvidence(ctx context.Context, input graphclient.CreateEvidenceInput, evidenceFiles []*graphql.Upload) (*graphclient.CreateEvidence, error)
	UpdateEvidence(ctx context.Context, id string, input graphclient.UpdateEvidenceInput, evidenceFiles []*graphql.Upload) (*graphclient.UpdateEvidence, error)
	CreateInternalPolicy(ctx context.Context, input graphclient.CreateInternalPolicyInput) (*graphclient.CreateInternalPolicy, error)
	UpdateInternalPolicy(ctx context.Context, id string, input graphclient.UpdateInternalPolicyInput) (*graphclient.UpdateInternalPolicy, error)
	CreateRisk(ctx context.Context, input graphclient.CreateRiskInput) (*graphclient.CreateRisk, error)
	UpdateRisk(ctx context.Context, id string, input graphclient.UpdateRiskInput) (*graphclient.UpdateRisk, error)
	CreateTask(ctx context.Context, input graphclient.CreateTaskInput) (*graphclient.CreateTask, error)
	UpdateTask(ctx context.Context, id string, input graphclient.UpdateTaskInput) (*graphclient.UpdateTask, error)
	DeleteControl(ctx context.Context, id string) (string, error)
	DeleteEvidence(ctx context.Context, id string) (string, error)
	DeleteInternalPolicy(ctx context.Context, id string) (string, error)
	DeleteRisk(ctx context.Context, id string) (string, error)
	DeleteTask(ctx context.Context, id string) (string, error)
}

// New constructs the official Openlane client. It does not call the API.
func New(cfg config.Config) (GraphAPI, error) {
	opts := []gqlclient.ClientOption{
		gqlclient.WithAPIToken(cfg.APIToken),
		gqlclient.WithBaseURL(cfg.BaseURL),
	}
	if cfg.OrganizationID != "" {
		opts = append(opts, gqlclient.WithInterceptors(gqlclient.WithOrganizationHeader(cfg.OrganizationID)))
	}

	c, err := gqlclient.New(opts...)
	if err != nil {
		return nil, fmt.Errorf("create openlane client: %w", RedactError(err))
	}
	timeout := httpTimeout
	if cfg.UploadTimeout > 0 {
		timeout = cfg.UploadTimeout
	}
	if req := c.HTTPSlingRequester(); req != nil {
		if hc := req.HTTPClient(); hc != nil {
			hc.Timeout = timeout
		}
	}
	return &api{c: c, uploadTimeout: timeout}, nil
}

type api struct {
	c             *gqlclient.Client
	uploadTimeout time.Duration
}

func (a *api) GetControls(ctx context.Context, first *int64, after *string, where *graphclient.ControlWhereInput) (*graphclient.GetControls, error) {
	return a.c.GetControls(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetControlByID(ctx context.Context, id string) (*graphclient.GetControlByID, error) {
	return a.c.GetControlByID(ctx, id)
}

func (a *api) GetPrograms(ctx context.Context, first *int64, after *string, where *graphclient.ProgramWhereInput) (*graphclient.GetPrograms, error) {
	return a.c.GetPrograms(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetProgramByID(ctx context.Context, id string) (*graphclient.GetProgramByID, error) {
	return a.c.GetProgramByID(ctx, id)
}

func (a *api) GetEvidences(ctx context.Context, first *int64, after *string, where *graphclient.EvidenceWhereInput) (*graphclient.GetEvidences, error) {
	return a.c.GetEvidences(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetEvidenceByID(ctx context.Context, id string) (*graphclient.GetEvidenceByID, error) {
	return a.c.GetEvidenceByID(ctx, id)
}

func (a *api) GetInternalPolicies(ctx context.Context, first *int64, after *string, where *graphclient.InternalPolicyWhereInput) (*graphclient.GetInternalPolicies, error) {
	return a.c.GetInternalPolicies(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetInternalPolicyByID(ctx context.Context, id string) (*graphclient.GetInternalPolicyByID, error) {
	return a.c.GetInternalPolicyByID(ctx, id)
}

func (a *api) GetRisks(ctx context.Context, first *int64, after *string, where *graphclient.RiskWhereInput) (*graphclient.GetRisks, error) {
	return a.c.GetRisks(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetRiskByID(ctx context.Context, id string) (*graphclient.GetRiskByID, error) {
	return a.c.GetRiskByID(ctx, id)
}

func (a *api) GetStandards(ctx context.Context, first *int64, after *string, where *graphclient.StandardWhereInput) (*graphclient.GetStandards, error) {
	return a.c.GetStandards(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetStandardByID(ctx context.Context, id string) (*graphclient.GetStandardByID, error) {
	return a.c.GetStandardByID(ctx, id)
}

func (a *api) GetTasks(ctx context.Context, first *int64, after *string, where *graphclient.TaskWhereInput) (*graphclient.GetTasks, error) {
	return a.c.GetTasks(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetTaskByID(ctx context.Context, id string) (*graphclient.GetTaskByID, error) {
	return a.c.GetTaskByID(ctx, id)
}

func (a *api) GetEntities(ctx context.Context, first *int64, after *string, where *graphclient.EntityWhereInput) (*graphclient.GetEntities, error) {
	return a.c.GetEntities(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetEntityByID(ctx context.Context, id string) (*graphclient.GetEntityByID, error) {
	return a.c.GetEntityByID(ctx, id)
}

func (a *api) GetAssets(ctx context.Context, first *int64, after *string, where *graphclient.AssetWhereInput) (*graphclient.GetAssets, error) {
	return a.c.GetAssets(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetAssetByID(ctx context.Context, id string) (*graphclient.GetAssetByID, error) {
	return a.c.GetAssetByID(ctx, id)
}

func (a *api) GetContacts(ctx context.Context, first *int64, after *string, where *graphclient.ContactWhereInput) (*graphclient.GetContacts, error) {
	return a.c.GetContacts(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetContactByID(ctx context.Context, id string) (*graphclient.GetContactByID, error) {
	return a.c.GetContactByID(ctx, id)
}

func (a *api) GetControlImplementations(ctx context.Context, first *int64, after *string, where *graphclient.ControlImplementationWhereInput) (*graphclient.GetControlImplementations, error) {
	return a.c.GetControlImplementations(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetControlImplementationByID(ctx context.Context, id string) (*graphclient.GetControlImplementationByID, error) {
	return a.c.GetControlImplementationByID(ctx, id)
}

func (a *api) GetAssessments(ctx context.Context, first *int64, after *string, where *graphclient.AssessmentWhereInput) (*graphclient.GetAssessments, error) {
	return a.c.GetAssessments(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetAssessmentByID(ctx context.Context, id string) (*graphclient.GetAssessmentByID, error) {
	return a.c.GetAssessmentByID(ctx, id)
}

func (a *api) GetAssessmentResponses(ctx context.Context, first *int64, after *string, where *graphclient.AssessmentResponseWhereInput) (*graphclient.GetAssessmentResponses, error) {
	return a.c.GetAssessmentResponses(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetFindings(ctx context.Context, first *int64, after *string, where *graphclient.FindingWhereInput) (*graphclient.GetFindings, error) {
	return a.c.GetFindings(ctx, first, nil, after, nil, where, nil)
}

func (a *api) GetFindingByID(ctx context.Context, id string) (*graphclient.GetFindingByID, error) {
	return a.c.GetFindingByID(ctx, id)
}

func (a *api) GetRemediations(ctx context.Context, first *int64, after *string, where *graphclient.RemediationWhereInput) (*graphclient.GetRemediations, error) {
	return a.c.GetRemediations(ctx, first, nil, after, nil, where, nil)
}

func (a *api) CreateControl(ctx context.Context, input graphclient.CreateControlInput) (*graphclient.CreateControl, error) {
	return a.c.CreateControl(ctx, input)
}

func (a *api) UpdateControl(ctx context.Context, id string, input graphclient.UpdateControlInput) (*graphclient.UpdateControl, error) {
	return a.c.UpdateControl(ctx, id, input)
}

func (a *api) CreateEvidence(ctx context.Context, input graphclient.CreateEvidenceInput, evidenceFiles []*graphql.Upload) (*graphclient.CreateEvidence, error) {
	return a.withUploadTimeout(ctx, len(evidenceFiles) > 0, func(ctx context.Context) (*graphclient.CreateEvidence, error) {
		return a.c.CreateEvidence(ctx, input, evidenceFiles)
	})
}

func (a *api) UpdateEvidence(ctx context.Context, id string, input graphclient.UpdateEvidenceInput, evidenceFiles []*graphql.Upload) (*graphclient.UpdateEvidence, error) {
	return a.withUploadTimeout(ctx, len(evidenceFiles) > 0, func(ctx context.Context) (*graphclient.UpdateEvidence, error) {
		return a.c.UpdateEvidence(ctx, id, input, evidenceFiles)
	})
}

func (a *api) withUploadTimeout[T any](ctx context.Context, uploading bool, fn func(context.Context) (T, error)) (T, error) {
	if !uploading || a.uploadTimeout <= httpTimeout {
		return fn(ctx)
	}
	uploadCtx, cancel := context.WithTimeout(ctx, a.uploadTimeout)
	defer cancel()
	return fn(uploadCtx)
}

func (a *api) CreateInternalPolicy(ctx context.Context, input graphclient.CreateInternalPolicyInput) (*graphclient.CreateInternalPolicy, error) {
	return a.c.CreateInternalPolicy(ctx, input)
}

func (a *api) UpdateInternalPolicy(ctx context.Context, id string, input graphclient.UpdateInternalPolicyInput) (*graphclient.UpdateInternalPolicy, error) {
	return a.c.UpdateInternalPolicy(ctx, id, input)
}

func (a *api) CreateRisk(ctx context.Context, input graphclient.CreateRiskInput) (*graphclient.CreateRisk, error) {
	return a.c.CreateRisk(ctx, input)
}

func (a *api) UpdateRisk(ctx context.Context, id string, input graphclient.UpdateRiskInput) (*graphclient.UpdateRisk, error) {
	return a.c.UpdateRisk(ctx, id, input)
}

func (a *api) CreateTask(ctx context.Context, input graphclient.CreateTaskInput) (*graphclient.CreateTask, error) {
	return a.c.CreateTask(ctx, input)
}

func (a *api) UpdateTask(ctx context.Context, id string, input graphclient.UpdateTaskInput) (*graphclient.UpdateTask, error) {
	return a.c.UpdateTask(ctx, id, input)
}

func (a *api) DeleteControl(ctx context.Context, id string) (string, error) {
	resp, err := a.c.DeleteControl(ctx, id)
	if err != nil {
		return "", err
	}
	return resp.DeleteControl.DeletedID, nil
}

func (a *api) DeleteEvidence(ctx context.Context, id string) (string, error) {
	resp, err := a.c.DeleteEvidence(ctx, id)
	if err != nil {
		return "", err
	}
	return resp.DeleteEvidence.DeletedID, nil
}

func (a *api) DeleteInternalPolicy(ctx context.Context, id string) (string, error) {
	resp, err := a.c.DeleteInternalPolicy(ctx, id)
	if err != nil {
		return "", err
	}
	return resp.DeleteInternalPolicy.DeletedID, nil
}

func (a *api) DeleteRisk(ctx context.Context, id string) (string, error) {
	resp, err := a.c.DeleteRisk(ctx, id)
	if err != nil {
		return "", err
	}
	return resp.DeleteRisk.DeletedID, nil
}

func (a *api) DeleteTask(ctx context.Context, id string) (string, error) {
	resp, err := a.c.DeleteTask(ctx, id)
	if err != nil {
		return "", err
	}
	return resp.DeleteTask.DeletedID, nil
}
