package openlane

import (
	"errors"
	"strings"
	"testing"
)

func TestRedactTokens(t *testing.T) {
	in := "auth failed for token tola_abcDEF123 and pat tolp_xyz-789"
	out := Redact(in)
	if strings.Contains(out, "abcDEF123") || strings.Contains(out, "xyz-789") {
		t.Fatalf("token leaked: %q", out)
	}
	if !strings.Contains(out, "tola_[redacted]") || !strings.Contains(out, "tolp_[redacted]") {
		t.Fatalf("expected redaction markers: %q", out)
	}
}

func TestRedactError(t *testing.T) {
	err := RedactError(errors.New("bearer tola_secretvalue rejected"))
	if strings.Contains(err.Error(), "secretvalue") {
		t.Fatalf("token leaked: %q", err.Error())
	}
}
