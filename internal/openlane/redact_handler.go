package openlane

import (
	"context"
	"log/slog"
)

// redactingHandler wraps a slog.Handler and redacts sensitive values from log messages.
type redactingHandler struct {
	inner slog.Handler
}

// NewRedactingHandler returns h with token, Authorization, and content_base64 redaction.
func NewRedactingHandler(h slog.Handler) slog.Handler {
	return &redactingHandler{inner: h}
}

func (h *redactingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.inner.Enabled(ctx, level)
}

func (h *redactingHandler) Handle(ctx context.Context, r slog.Record) error {
	r2 := slog.NewRecord(r.Time, r.Level, RedactLogMessage(r.Message), r.PC)
	r.Attrs(func(a slog.Attr) bool {
		r2.AddAttrs(redactAttr(a))
		return true
	})
	return h.inner.Handle(ctx, r2)
}

func (h *redactingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	out := make([]slog.Attr, 0, len(attrs))
	for _, a := range attrs {
		out = append(out, redactAttr(a))
	}
	return &redactingHandler{inner: h.inner.WithAttrs(out)}
}

func (h *redactingHandler) WithGroup(name string) slog.Handler {
	return &redactingHandler{inner: h.inner.WithGroup(name)}
}

func redactAttr(a slog.Attr) slog.Attr {
	switch a.Value.Kind() {
	case slog.KindString:
		return slog.String(a.Key, RedactLogMessage(a.Value.String()))
	case slog.KindAny:
		if s, ok := a.Value.Any().(string); ok {
			return slog.String(a.Key, RedactLogMessage(s))
		}
	}
	// Never log raw token env keys.
	if isSensitiveAttrKey(a.Key) {
		return slog.String(a.Key, "[redacted]")
	}
	return a
}

func isSensitiveAttrKey(key string) bool {
	switch key {
	case "OPENLANE_API_TOKEN", "api_token", "token", "authorization", "Authorization":
		return true
	default:
		return false
	}
}
