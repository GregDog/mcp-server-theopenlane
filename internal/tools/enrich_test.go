package tools

import (
	"testing"

	"github.com/theopenlane/go-client/graphclient"
)

func TestBuildPolicyWhere(t *testing.T) {
	where := buildPolicyWhere(policyListInput{Status: "NEEDS_APPROVAL"})
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.Status == nil || string(*where.Status) != "NEEDS_APPROVAL" {
		t.Fatalf("status filter: %+v", where.Status)
	}
	if buildPolicyWhere(policyListInput{}) != nil {
		t.Fatal("expected nil where for empty status")
	}
}

func TestBuildFindingWhere(t *testing.T) {
	open := true
	where := buildFindingWhere(findingListInput{
		ProgramID:    "prog-1",
		AssessmentID: "asm-1",
		Open:         &open,
		Status:       "open",
		Severity:     "high",
	})
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.AssessmentID == nil || *where.AssessmentID != "asm-1" {
		t.Fatalf("assessment filter: %+v", where.AssessmentID)
	}
	if where.Open == nil || !*where.Open {
		t.Fatalf("open filter: %+v", where.Open)
	}
	if where.Severity == nil || *where.Severity != "high" {
		t.Fatalf("severity filter: %+v", where.Severity)
	}
	if len(where.HasProgramsWith) != 1 || where.HasProgramsWith[0].ID == nil || *where.HasProgramsWith[0].ID != "prog-1" {
		t.Fatalf("program filter: %+v", where.HasProgramsWith)
	}
}

func TestBuildEntityWhere(t *testing.T) {
	approved := true
	where := buildEntityWhere(entityListInput{
		RiskRating:          "high",
		QuestionnaireStatus: "complete",
		ApprovedForUse:      &approved,
	})
	if where == nil {
		t.Fatal("expected where input")
	}
	if where.RiskRating == nil || *where.RiskRating != "high" {
		t.Fatalf("risk rating: %+v", where.RiskRating)
	}
	if where.ApprovedForUse == nil || !*where.ApprovedForUse {
		t.Fatalf("approved: %+v", where.ApprovedForUse)
	}
}

func TestMapGetEntityVendorFields(t *testing.T) {
	hasSoc2 := true
	sso := false
	mfa := true
	approved := true
	spend := 120000.0
	name := "Acme"
	item := mapGetEntity(graphclient.GetEntityByID_Entity{
		ID:             "ent-1",
		Name:           &name,
		HasSoc2:        &hasSoc2,
		SsoEnforced:    &sso,
		MfaEnforced:    &mfa,
		ApprovedForUse: &approved,
		AnnualSpend:    &spend,
		RiskRating:     strPtr("high"),
	})
	if item.Name != "Acme" || item.HasSoc2 == nil || !*item.HasSoc2 {
		t.Fatalf("unexpected entity mapping: %+v", item)
	}
	if item.SsoEnforced == nil || *item.SsoEnforced {
		t.Fatalf("sso: %+v", item.SsoEnforced)
	}
	if item.AnnualSpend == nil || *item.AnnualSpend != spend {
		t.Fatalf("spend: %+v", item.AnnualSpend)
	}
}

func strPtr(s string) *string { return &s }
