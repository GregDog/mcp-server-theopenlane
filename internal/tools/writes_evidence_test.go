package tools

import (
	"context"
	"strings"
	"testing"

	"github.com/theopenlane/core/common/enums"
	"github.com/theopenlane/go-client/graphclient"
)

func TestCreateEvidenceWithControlIDs(t *testing.T) {
	owner := "org_1"
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      "ctrl_awareness",
				OwnerID: &owner,
			},
		},
	}
	h := &handlers{api: api, allowWrite: true}
	_, item, err := h.createEvidence(context.Background(), nil, createEvidenceInput{
		Name:       "KnowBe4 training completion",
		ControlIDs: []string{"ctrl_awareness"},
	})
	if err != nil {
		t.Fatalf("create evidence: %v", err)
	}
	if item.ID != "evidence_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
	if len(api.lastCreateEvidenceInput.ControlIDs) != 1 || api.lastCreateEvidenceInput.ControlIDs[0] != "ctrl_awareness" {
		t.Fatalf("unexpected create input control IDs: %+v", api.lastCreateEvidenceInput.ControlIDs)
	}
}

func TestUpdateEvidenceAllowsAddControlIDsOnly(t *testing.T) {
	owner := "org_1"
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      "ctrl_awareness",
				OwnerID: &owner,
			},
		},
	}
	h := &handlers{api: api, allowWrite: true}
	_, item, err := h.updateEvidence(context.Background(), nil, updateEvidenceInput{
		ID:            "evidence_1",
		AddControlIDs: []string{"ctrl_awareness"},
	})
	if err != nil {
		t.Fatalf("update evidence: %v", err)
	}
	if item.ID != "evidence_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
	if len(api.lastUpdateEvidenceInput.AddControlIDs) != 1 || api.lastUpdateEvidenceInput.AddControlIDs[0] != "ctrl_awareness" {
		t.Fatalf("unexpected update add control IDs: %+v", api.lastUpdateEvidenceInput.AddControlIDs)
	}
}

func TestCreateEvidenceRejectsCatalogControl(t *testing.T) {
	title := "Security awareness education is an ongoing activity."
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      "ctrl_catalog",
				RefCode: "12.6.2",
				Title:   &title,
			},
		},
	}
	h := &handlers{api: api, allowWrite: true}
	_, _, err := h.createEvidence(context.Background(), nil, createEvidenceInput{
		Name:       "KnowBe4 training completion",
		ControlIDs: []string{"ctrl_catalog"},
	})
	if err == nil {
		t.Fatal("expected error for catalog control")
	}
	if !strings.Contains(err.Error(), "catalog/system control") {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestCreateEvidenceAllowsOrgOwnedFrameworkControl(t *testing.T) {
	framework := enums.ControlSourceFramework
	owner := "org_1"
	api := &fakeAPI{
		control: &graphclient.GetControlByID{
			Control: graphclient.GetControlByID_Control{
				ID:      "ctrl_program",
				OwnerID: &owner,
				Source:  &framework,
			},
		},
	}
	h := &handlers{api: api, allowWrite: true}
	_, _, err := h.createEvidence(context.Background(), nil, createEvidenceInput{
		Name:       "KnowBe4 training completion",
		ControlIDs: []string{"ctrl_program"},
	})
	if err != nil {
		t.Fatalf("create evidence: %v", err)
	}
}

func TestUpdateEvidenceAllowsRemoveControlIDsOnly(t *testing.T) {
	api := &fakeAPI{}
	h := &handlers{api: api, allowWrite: true}
	_, item, err := h.updateEvidence(context.Background(), nil, updateEvidenceInput{
		ID:               "evidence_1",
		RemoveControlIDs: []string{"ctrl_awareness"},
	})
	if err != nil {
		t.Fatalf("update evidence: %v", err)
	}
	if item.ID != "evidence_1" {
		t.Fatalf("unexpected item: %+v", item)
	}
	if len(api.lastUpdateEvidenceInput.RemoveControlIDs) != 1 || api.lastUpdateEvidenceInput.RemoveControlIDs[0] != "ctrl_awareness" {
		t.Fatalf("unexpected update remove control IDs: %+v", api.lastUpdateEvidenceInput.RemoveControlIDs)
	}
}
