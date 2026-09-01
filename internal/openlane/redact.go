package openlane

import (
	"regexp"
	"strings"
)

var tokenPattern = regexp.MustCompile(`(?i)\b(tola_|tolp_)[A-Za-z0-9_\-\.]+`)

// Redact removes Openlane API tokens and PATs from a string.
func Redact(s string) string {
	if s == "" {
		return s
	}
	return tokenPattern.ReplaceAllStringFunc(s, func(match string) string {
		lower := strings.ToLower(match)
		switch {
		case strings.HasPrefix(lower, "tola_"):
			return "tola_[redacted]"
		case strings.HasPrefix(lower, "tolp_"):
			return "tolp_[redacted]"
		default:
			return "[redacted]"
		}
	})
}

// RedactError returns err with token values removed from the message.
func RedactError(err error) error {
	if err == nil {
		return nil
	}
	msg := Redact(err.Error())
	if msg == err.Error() {
		return err
	}
	return redactedError{msg: msg, unwrap: err}
}

type redactedError struct {
	msg    string
	unwrap error
}

func (e redactedError) Error() string { return e.msg }
func (e redactedError) Unwrap() error { return e.unwrap }
