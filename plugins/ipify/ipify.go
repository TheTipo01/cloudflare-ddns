package main

import (
	"io"
	"net/http"
)

// return public IP address of the machine from the specified endpoint
func GetIPs(loadV4, loadV6 bool) (string, string, error) {
	ipv4, ipv6 := "", ""
	if loadV4 {
		var err error
		ipv4, err = getIP("https://api4.ipify.org")
		if err != nil {
			return "", "", err
		}
	}
	if loadV6 {
		var err error
		ipv6, err = getIP("https://api6.ipify.org")
		if err != nil {
			return "", "", err
		}
	}
	return ipv4, ipv6, nil
}

// return public IP address of the machine from the specified endpoint
func getIP(endpoint string) (string, error) {
	response, err := http.Get(endpoint)
	if err == nil {
		defer response.Body.Close()
		out, err := io.ReadAll(response.Body)
		if err != nil {
			return "", err
		}
		return string(out), nil
	} else {
		return "", err
	}
}
