package main

import (
	"io"
	"net/http"
	"sync"
)

// return public IP address of the machine from the specified endpoint
func GetIPs() (string, string, error) {
	wg := sync.WaitGroup{}
	wg.Add(2)
	errMut := sync.Mutex{}
	var ipv4, ipv6 string
	var err error
	go func() {
		defer wg.Done()
		ip, e := getIP("https://api4.ipify.org")
		if e != nil {
			errMut.Lock()
			defer errMut.Unlock()
			err = e
		} else {
			ipv4 = ip
		}
	}()
	go func() {
		defer wg.Done()
		ip, e := getIP("https://api6.ipify.org")
		if e != nil {
			errMut.Lock()
			defer errMut.Unlock()
			err = e
		} else {
			ipv6 = ip
		}
	}()
	wg.Wait()

	if err != nil {
		return "", "", err
	} else {
		return ipv4, ipv6, nil
	}
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
