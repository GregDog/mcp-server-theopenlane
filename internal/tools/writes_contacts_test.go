package tools

import (
	"context"
	"testing"

	"github.com/GregDog/mcp-server-theopenlane/internal/openlane"
)

func TestCreateContactRequiresIdentifier(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, _, err := h.createContact(context.Background(), nil, createContactInput{})
	if err != errContactIdentifierRequired {
		t.Fatalf("create contact: got %v", err)
	}
}

func TestCreateContactSuccess(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, item, err := h.createContact(context.Background(), nil, createContactInput{
		FullName:  "David Davidson",
		Email:     "ddavidson@saepio.co.uk",
		EntityIDs: []string{"ent_saepio"},
	})
	if err != nil {
		t.Fatalf("create contact: %v", err)
	}
	if item.ID != "contact_1" || item.FullName != "David Davidson" || item.Email != "ddavidson@saepio.co.uk" {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestCreateRiskWithEntityIDs(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}
	_, item, err := h.createRisk(context.Background(), nil, createRiskInput{
		Name:      "Saepio vendor risk assessment",
		EntityIDs: []string{"ent_saepio"},
	})
	if err != nil {
		t.Fatalf("create risk: %v", err)
	}
	if item.ID != "risk_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
}

func TestUpdateEntityAllowsAddContactIDsOnly(t *testing.T) {
	h := &handlers{
		api:        &fakeAPI{entity: &openlane.EntityDetail{ID: "ent_1"}},
		allowWrite: true,
	}
	_, item, err := h.updateEntity(context.Background(), nil, updateEntityInput{
		ID:            "ent_1",
		AddContactIDs: []string{"contact_1"},
	})
	if err != nil {
		t.Fatalf("update entity contacts: %v", err)
	}
	if item.ID != "ent_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
}
