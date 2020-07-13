package main

import (
	"net"
	"strings"
)

// IsIPv4 checks if IP a valid IPv4.
func IsIPv4(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if len(ip.To4()) == net.IPv4len && strings.Contains(ip.String(), ".") {
		return true
	}

	return false
}

// IsIPv6 checks if IP a valid IPv6.
func IsIPv6(ip net.IP) bool {
	if ip == nil {
		return false
	}

	if len(ip.To16()) == net.IPv6len && strings.Contains(ip.String(), ":") {
		return true
	}

	return false
}
