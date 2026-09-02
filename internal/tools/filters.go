package tools

import (
	"strings"
	"time"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/core/common/models"
	"github.com/theopenlane/go-client/graphclient"
)

func buildEntityWhere(in entityListInput) *graphclient.EntityWhereInput {
	var w graphclient.EntityWhereInput
	has := false
	if s := strings.TrimSpace(in.RiskRating); s != "" {
		w.RiskRating = &s
		has = true
	}
	if s := strings.TrimSpace(in.Tier); s != "" {
		t := enums.VendorTier(s)
		w.Tier = &t
		has = true
	}
	if in.ApprovedForUse != nil {
		w.ApprovedForUse = in.ApprovedForUse
		has = true
	}
	if s := strings.TrimSpace(in.QuestionnaireStatus); s != "" {
		w.EntitySecurityQuestionnaireStatusName = &s
		has = true
	}
	if s := strings.TrimSpace(in.NextReviewBefore); s != "" {
		if dt, err := parseDateTime(s); err == nil {
			w.NextReviewAtLte = dt
			has = true
		}
	}
	if in.HasSoc2 != nil {
		w.HasSoc2 = in.HasSoc2
		has = true
	}
	if in.SsoEnforced != nil {
		w.SsoEnforced = in.SsoEnforced
		has = true
	}
	if in.MfaEnforced != nil {
		w.MfaEnforced = in.MfaEnforced
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func buildRiskWhere(in riskListInput) *graphclient.RiskWhereInput {
	var w graphclient.RiskWhereInput
	has := false
	if s := strings.TrimSpace(in.ProgramID); s != "" {
		w.HasProgramsWith = []*graphclient.ProgramWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.EntityID); s != "" {
		w.HasEntitiesWith = []*graphclient.EntityWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.ControlID); s != "" {
		w.HasControlsWith = []*graphclient.ControlWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.Status); s != "" {
		st := enums.RiskStatus(s)
		w.Status = &st
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func buildFindingWhere(in findingListInput) *graphclient.FindingWhereInput {
	var w graphclient.FindingWhereInput
	has := false
	if s := strings.TrimSpace(in.ProgramID); s != "" {
		w.HasProgramsWith = []*graphclient.ProgramWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.AssessmentID); s != "" {
		w.AssessmentID = &s
		has = true
	}
	if in.Open != nil {
		w.Open = in.Open
		has = true
	}
	if s := strings.TrimSpace(in.Status); s != "" {
		w.FindingStatusName = &s
		has = true
	}
	if s := strings.TrimSpace(in.Severity); s != "" {
		w.Severity = &s
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func buildImplementationWhere(in implementationListInput) *graphclient.ControlImplementationWhereInput {
	var w graphclient.ControlImplementationWhereInput
	has := false
	if s := strings.TrimSpace(in.ControlID); s != "" {
		w.HasControlsWith = []*graphclient.ControlWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.Status); s != "" {
		st := enums.DocumentStatus(s)
		w.Status = &st
		has = true
	}
	if in.Verified != nil {
		w.Verified = in.Verified
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func buildAssessmentWhere(in assessmentListInput) *graphclient.AssessmentWhereInput {
	s := strings.TrimSpace(in.AssessmentType)
	if s == "" {
		return nil
	}
	t := enums.AssessmentType(s)
	return &graphclient.AssessmentWhereInput{AssessmentType: &t}
}

func buildProgramWhere(in programListInput) *graphclient.ProgramWhereInput {
	s := strings.TrimSpace(in.Name)
	if s == "" {
		return nil
	}
	return &graphclient.ProgramWhereInput{NameContainsFold: &s}
}

func buildEvidenceWhere(in evidenceListInput) *graphclient.EvidenceWhereInput {
	var w graphclient.EvidenceWhereInput
	has := false
	if s := strings.TrimSpace(in.ProgramID); s != "" {
		w.HasProgramsWith = []*graphclient.ProgramWhereInput{{ID: &s}}
		has = true
	}
	if s := strings.TrimSpace(in.ControlID); s != "" {
		w.HasControlsWith = []*graphclient.ControlWhereInput{{ID: &s}}
		has = true
	}
	if !has {
		return nil
	}
	return &w
}

func parseDateTime(s string) (*models.DateTime, error) {
	if dt, err := models.ToDateTime(s); err == nil {
		return dt, nil
	}
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return nil, err
	}
	dt := models.DateTime(t)
	return &dt, nil
}
