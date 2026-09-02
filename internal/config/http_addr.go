package config

import (
	"net"
	"strings"
)

// IsLoopbackHTTPAddr reports whether addr listens only on loopback interfaces.
// Addresses without a host (for example ":8090") bind all interfaces and return false.
func IsLoopbackHTTPAddr(addr string) bool {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return false
	}
	host = strings.TrimSpace(host)
	if host == "" {
		return false
	}
	if strings.EqualFold(host, "localhost") {
		return true
	}
	ip := net.ParseIP(host)
	return ip != nil && ip.IsLoopback()
}
