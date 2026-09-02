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
	if err == nil {
		t.Fatal("expected error")
	}
	if strings.Contains(err.Error(), "tola_secretvalue") {
		t.Fatalf("token not redacted: %q", err.Error())
	}
}

func TestRedactLogMessage(t *testing.T) {
	in := `Authorization: Bearer tolp_abc123 and "content_base64":"aGVsbG8="`
	out := RedactLogMessage(in)
	if strings.Contains(out, "tolp_abc123") || strings.Contains(out, "aGVsbG8=") {
		t.Fatalf("expected redaction: %q", out)
	}
	if !strings.Contains(out, "[redacted]") {
		t.Fatalf("expected redaction markers: %q", out)
	}
}
