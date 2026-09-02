package tools

import (
	"context"
	"testing"
)

func TestWriteValidationRequiresFields(t *testing.T) {
	h := &handlers{api: &fakeAPI{}, allowWrite: true}

	_, _, err := h.createControl(context.Background(), nil, createControlInput{})
	if err != errRefCodeRequired {
		t.Fatalf("create control: got %v", err)
	}

	_, _, err = h.createEvidence(context.Background(), nil, createEvidenceInput{})
	if err != errNameRequired {
		t.Fatalf("create evidence: got %v", err)
	}

	_, _, err = h.createPolicy(context.Background(), nil, createPolicyInput{})
	if err != errNameRequired {
		t.Fatalf("create policy: got %v", err)
	}

	_, _, err = h.createRisk(context.Background(), nil, createRiskInput{})
	if err != errNameRequired {
		t.Fatalf("create risk: got %v", err)
	}

	_, _, err = h.createTask(context.Background(), nil, createTaskInput{})
	if err != errTitleRequired {
		t.Fatalf("create task: got %v", err)
	}

	_, _, err = h.updateControl(context.Background(), nil, updateControlInput{ID: "ctrl_1"})
	if err != errUpdateFieldsRequired {
		t.Fatalf("update control: got %v", err)
	}
}
