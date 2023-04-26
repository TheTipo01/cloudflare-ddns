package main

import (
	"plugin"
	"sync"

	"github.com/nylone/cloudflare-ddns/src/cloudflare"
)

func main() {
	p, err := plugin.Open("plugins/ipify/ipify.so")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("GetIPs")
	if err != nil {
		panic(err)
	}
	getIPs := f.(func() (string, string, error))

	cloudflare.LoadMappings()
	ipv4, ipv6, err := getIPs()
	if err != nil {
		wg := sync.WaitGroup{}
		defer wg.Wait()
		cloudflare.PatchARecords(ipv4, &wg)
		cloudflare.PatchAAAARecords(ipv6, &wg)
	}
}
