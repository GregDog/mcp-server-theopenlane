package openlane

import "testing"

func TestControlLinkableToEvidence(t *testing.T) {
	org := "01KTY1J6M1WTQAPYC2T2MKREKD"
	empty := ""
	if !ControlLinkableToEvidence(&org) {
		t.Fatal("expected org-owned control to be linkable")
	}
	if ControlLinkableToEvidence(nil) {
		t.Fatal("expected nil owner not linkable")
	}
	if ControlLinkableToEvidence(&empty) {
		t.Fatal("expected empty owner not linkable")
	}
}
