package openlane

import (
	"context"
	"fmt"
	"time"

	gqlclient "github.com/theopenlane/go-client"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/config"
)

const httpTimeout = 30 * time.Second

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
	if req := c.HTTPSlingRequester(); req != nil {
		if hc := req.HTTPClient(); hc != nil {
			hc.Timeout = httpTimeout
		}
	}
	return &api{c: c}, nil
}

type api struct {
	c *gqlclient.Client
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
