package openlane

import (
	"fmt"
	"reflect"
	"time"
)

func Deref[T any](p *T) T {
	var zero T
	if p == nil {
		return zero
	}
	return *p
}

func Format(v any) string {
	if v == nil {
		return ""
	}
	rv := reflect.ValueOf(v)
	for rv.Kind() == reflect.Pointer {
		if rv.IsNil() {
			return ""
		}
		rv = rv.Elem()
	}
	v = rv.Interface()
	switch t := v.(type) {
	case string:
		return t
	case time.Time:
		if t.IsZero() {
			return ""
		}
		return t.UTC().Format(time.RFC3339)
	default:
		return fmt.Sprint(v)
	}
}
