package main

import "testing"

func TestGetIPs(t *testing.T) {
	ipv4, ipv6, err := GetIPs(true, true)
	if err != nil {
		t.Error(err)
	} else {
		t.Logf("IPv4: %s, IPv6: %s", ipv4, ipv6)
	}
}
