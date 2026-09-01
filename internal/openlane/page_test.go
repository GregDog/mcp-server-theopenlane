package openlane

import "testing"

func TestClampLimit(t *testing.T) {
	cases := []struct {
		in   int
		want int64
	}{
		{0, DefaultPageSize},
		{-1, DefaultPageSize},
		{1, 1},
		{20, 20},
		{50, 50},
		{51, MaxPageSize},
		{1000, MaxPageSize},
	}
	for _, tc := range cases {
		if got := ClampLimit(tc.in); got != tc.want {
			t.Fatalf("ClampLimit(%d)=%d want %d", tc.in, got, tc.want)
		}
	}
}

func TestCursorPtr(t *testing.T) {
	if CursorPtr("") != nil {
		t.Fatal("empty cursor should be nil")
	}
	p := CursorPtr("abc")
	if p == nil || *p != "abc" {
		t.Fatalf("got %#v", p)
	}
}
