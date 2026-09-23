package tools

import (
	"context"
	"errors"

	"github.com/99designs/gqlgen/graphql"
	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func (f *fakeAPI) CreateControl(context.Context, graphclient.CreateControlInput) (*graphclient.CreateControl, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateControl(_ context.Context, id string, input graphclient.UpdateControlInput) (*graphclient.UpdateControl, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastUpdateControlInput = input
	return &graphclient.UpdateControl{
		UpdateControl: graphclient.UpdateControl_UpdateControl{
			Control: graphclient.UpdateControl_UpdateControl_Control{
				ID:             id,
				RefCode:        "AC-1",
				ControlOwnerID: input.ControlOwnerID,
				DelegateID:     input.DelegateID,
			},
		},
	}, nil
}
func (f *fakeAPI) CreateEvidence(_ context.Context, input graphclient.CreateEvidenceInput, _ []*graphql.Upload) (*graphclient.CreateEvidence, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastCreateEvidenceInput = input
	return &graphclient.CreateEvidence{
		CreateEvidence: graphclient.CreateEvidence_CreateEvidence{
			Evidence: graphclient.CreateEvidence_CreateEvidence_Evidence{
				ID:   "evidence_1",
				Name: input.Name,
			},
		},
	}, nil
}
func (f *fakeAPI) UpdateEvidence(_ context.Context, id string, input graphclient.UpdateEvidenceInput, _ []*graphql.Upload) (*graphclient.UpdateEvidence, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastUpdateEvidenceInput = input
	name := "updated"
	if input.Name != nil {
		name = *input.Name
	}
	return &graphclient.UpdateEvidence{
		UpdateEvidence: graphclient.UpdateEvidence_UpdateEvidence{
			Evidence: graphclient.UpdateEvidence_UpdateEvidence_Evidence{
				ID:   id,
				Name: name,
			},
		},
	}, nil
}
func (f *fakeAPI) CreateInternalPolicy(context.Context, graphclient.CreateInternalPolicyInput) (*graphclient.CreateInternalPolicy, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateInternalPolicy(_ context.Context, id string, input graphclient.UpdateInternalPolicyInput) (*graphclient.UpdateInternalPolicy, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastUpdatePolicyInput = input
	if f.policy != nil && input.Status != nil {
		f.policy.InternalPolicy.Status = input.Status
	}
	name := "updated"
	if input.Name != nil {
		name = *input.Name
	}
	return &graphclient.UpdateInternalPolicy{
		UpdateInternalPolicy: graphclient.UpdateInternalPolicy_UpdateInternalPolicy{
			InternalPolicy: graphclient.UpdateInternalPolicy_UpdateInternalPolicy_InternalPolicy{
				ID:         id,
				Name:       name,
				ApproverID: input.ApproverID,
				DelegateID: input.DelegateID,
			},
		},
	}, nil
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
func (f *fakeAPI) UpdateRisk(_ context.Context, id string, input graphclient.UpdateRiskInput) (*graphclient.UpdateRisk, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastUpdateRiskInput = input
	name := "updated"
	if input.Name != nil {
		name = *input.Name
	}
	return &graphclient.UpdateRisk{
		UpdateRisk: graphclient.UpdateRisk_UpdateRisk{
			Risk: graphclient.UpdateRisk_UpdateRisk_Risk{
				ID:            id,
				Name:          name,
				StakeholderID: input.StakeholderID,
				DelegateID:    input.DelegateID,
			},
		},
	}, nil
}
func (f *fakeAPI) CreateTask(context.Context, graphclient.CreateTaskInput) (*graphclient.CreateTask, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) UpdateTask(context.Context, string, graphclient.UpdateTaskInput) (*graphclient.UpdateTask, error) {
	return nil, errors.New("unused")
}
func (f *fakeAPI) GetMappedControls(_ context.Context, _ *int64, _ *string, where *graphclient.MappedControlWhereInput) (*graphclient.GetMappedControls, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.mappedControls != nil {
		return f.mappedControls, nil
	}
	return &graphclient.GetMappedControls{
		MappedControls: graphclient.GetMappedControls_MappedControls{
			Edges: []*graphclient.GetMappedControls_MappedControls_Edges{},
		},
	}, nil
}
func (f *fakeAPI) GetMappedControlByID(_ context.Context, id string) (*graphclient.GetMappedControlByID, error) {
	if f.err != nil {
		return nil, f.err
	}
	if f.mappedControl != nil {
		return f.mappedControl, nil
	}
	return nil, errors.New("unused")
}
func (f *fakeAPI) CreateMappedControl(_ context.Context, input graphclient.CreateMappedControlInput) (*graphclient.CreateMappedControl, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastCreateMappedControlInput = input
	from := input.FromControlIDs
	to := input.ToControlIDs
	edgesFrom := make([]*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromControls_Edges, 0, len(from))
	for _, id := range from {
		edgesFrom = append(edgesFrom, &graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromControls_Edges{
			Node: &graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromControls_Edges_Node{ID: id},
		})
	}
	edgesTo := make([]*graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToControls_Edges, 0, len(to))
	for _, id := range to {
		edgesTo = append(edgesTo, &graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToControls_Edges{
			Node: &graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToControls_Edges_Node{ID: id},
		})
	}
	mt := enums.MappingTypeEqual
	if input.MappingType != nil {
		mt = *input.MappingType
	}
	return &graphclient.CreateMappedControl{
		CreateMappedControl: graphclient.CreateMappedControl_CreateMappedControl{
			MappedControl: graphclient.CreateMappedControl_CreateMappedControl_MappedControl{
				ID:          "map_1",
				MappingType: mt,
				FromControls: graphclient.CreateMappedControl_CreateMappedControl_MappedControl_FromControls{
					Edges: edgesFrom,
				},
				ToControls: graphclient.CreateMappedControl_CreateMappedControl_MappedControl_ToControls{
					Edges: edgesTo,
				},
			},
		},
	}, nil
}
func (f *fakeAPI) UpdateMappedControl(_ context.Context, id string, input graphclient.UpdateMappedControlInput) (*graphclient.UpdateMappedControl, error) {
	if f.err != nil {
		return nil, f.err
	}
	f.lastUpdateMappedControlInput = input
	mt := enums.MappingTypeEqual
	if input.MappingType != nil {
		mt = *input.MappingType
	}
	return &graphclient.UpdateMappedControl{
		UpdateMappedControl: graphclient.UpdateMappedControl_UpdateMappedControl{
			MappedControl: graphclient.UpdateMappedControl_UpdateMappedControl_MappedControl{
				ID:          id,
				MappingType: mt,
			},
		},
	}, nil
}
func (f *fakeAPI) DeleteMappedControl(_ context.Context, id string) (string, error) {
	if f.err != nil {
		return "", f.err
	}
	f.deletedID = id
	return id, nil
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
func (f *fakeAPI) CreateReview(_ context.Context, input graphclient.CreateReviewInput) (*graphclient.CreateReview, error) {
	if f.err != nil {
		return nil, f.err
	}
	return &graphclient.CreateReview{
		CreateReview: graphclient.CreateReview_CreateReview{
			Review: graphclient.CreateReview_CreateReview_Review{
				ID:      "review_1",
				Title:   input.Title,
				Details: input.Details,
				Summary: input.Summary,
			},
		},
	}, nil
}
func (f *fakeAPI) UpdateReview(_ context.Context, id string, input graphclient.UpdateReviewInput) (*graphclient.UpdateReview, error) {
	if f.err != nil {
		return nil, f.err
	}
	title := "updated"
	if input.Title != nil {
		title = *input.Title
	}
	return &graphclient.UpdateReview{
		UpdateReview: graphclient.UpdateReview_UpdateReview{
			Review: graphclient.UpdateReview_UpdateReview_Review{
				ID:      id,
				Title:   title,
				Details: input.Details,
				Summary: input.Summary,
			},
		},
	}, nil
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
