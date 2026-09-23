package tools

import (
	"context"
	"testing"
)

func TestCreateVendorRiskReviewRequiresTitleAndEntity(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, _, err := h.createVendorRiskReview(context.Background(), nil, createVendorRiskReviewInput{})
	if err != errTitleRequired {
		t.Fatalf("create review missing title: got %v", err)
	}
	_, _, err = h.createVendorRiskReview(context.Background(), nil, createVendorRiskReviewInput{
		Title: "Saepio vendor risk assessment",
	})
	if err != errEntityIDRequired {
		t.Fatalf("create review missing entity: got %v", err)
	}
}

func TestCreateVendorRiskReviewSuccess(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	body := "Assessment body for vendor onboarding."
	_, item, err := h.createVendorRiskReview(context.Background(), nil, createVendorRiskReviewInput{
		Title:     "Example vendor risk assessment",
		Body:      body,
		EntityIDs: []string{"ent_saepio"},
	})
	if err != nil {
		t.Fatalf("create review: %v", err)
	}
	if item.ID != "review_1" || item.Title == "" || item.Details != body {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestUpdateRiskAllowsRemoveEntityIDsOnly(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, item, err := h.updateRisk(context.Background(), nil, updateRiskInput{
		ID:              "risk_1",
		RemoveEntityIDs: []string{"ent_saepio"},
	})
	if err != nil {
		t.Fatalf("update risk unlink: %v", err)
	}
	if item.ID != "risk_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
}
