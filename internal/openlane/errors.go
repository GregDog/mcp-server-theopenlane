package openlane

import (
	"context"
	"errors"
	"fmt"
)

// APIError maps Openlane/client failures into a safe, redacted error.
func APIError(err error) error {
	if err == nil {
		return nil
	}
	switch {
	case errors.Is(err, context.Canceled):
		return fmt.Errorf("request canceled")
	case errors.Is(err, context.DeadlineExceeded):
		return fmt.Errorf("openlane request timed out")
	default:
		return fmt.Errorf("openlane api: %s", Redact(err.Error()))
	}
}
