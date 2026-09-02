package config

import (
	"strings"
	"testing"
)

func TestDefaultHTTPAddrIsLoopback(t *testing.T) {
	if !IsLoopbackHTTPAddr(DefaultHTTPAddr) {
		t.Fatalf("DefaultHTTPAddr %q must bind loopback only", DefaultHTTPAddr)
	}
	if strings.Contains(DefaultHTTPAddr, "0.0.0.0") {
		t.Fatalf("DefaultHTTPAddr must not contain 0.0.0.0: %q", DefaultHTTPAddr)
	}
}

func TestIsLoopbackHTTPAddr(t *testing.T) {
	tests := []struct {
		addr string
		want bool
	}{
		{DefaultHTTPAddr, true},
		{"127.0.0.1:8090", true},
		{"localhost:8090", true},
		{"[::1]:8090", true},
		{":8090", false},
		{"0.0.0.0:8090", false},
		{"10.0.0.5:8090", false},
	}
	for _, tc := range tests {
		if got := IsLoopbackHTTPAddr(tc.addr); got != tc.want {
			t.Fatalf("IsLoopbackHTTPAddr(%q) = %v, want %v", tc.addr, got, tc.want)
		}
	}
}
