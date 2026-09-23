package openlane

// ControlLinkableToEvidence reports whether evidence may be linked to the control.
// Openlane blocks system-owned catalog controls as evidence edges (theopenlane/core#1647).
// Program-imported controls keep source=FRAMEWORK but have ownerID set and are linkable.
func ControlLinkableToEvidence(ownerID *string) bool {
	return ownerID != nil && *ownerID != ""
}
