package main

import (
	"plugin"
	"time"

	"github.com/kkyr/fig"
	"github.com/nylone/cloudflare-ddns/cloudflare"
)

const (
	configFile = "config/config.yml"
)

type config struct {
	Backend string `fig:"backend" validate:"required"`
	Timeout int    `fig:"timeout" default:"30"`
	DoV4    bool   `fig:"doV4"`
	DoV6    bool   `fig:"doV6"`
}

var (
	getIPs func(bool, bool) (string, string, error)
	cfg    config
)

func init() {
	// load config file
	err := fig.Load(&cfg, fig.File(configFile))
	if err != nil {
		panic(err)
	}

	// load plugin
	p, err := plugin.Open(cfg.Backend + ".so")
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("GetIPs")
	if err != nil {
		panic(err)
	}
	getIPs = f.(func(bool, bool) (string, string, error))
}

func main() {
	// load cloudflare mappings
	err := cloudflare.LoadMappings()
	if err != nil {
		panic(err)
	}

	ipv4, ipv6 := "", ""
	var newIpv4, newIpv6 string
	for {
		newIpv4, newIpv6, err = getIPs(cfg.DoV4, cfg.DoV6)
		if err != nil {
			panic(err)
		}
		if newIpv4 != ipv4 {
			ipv4 = newIpv4
			cloudflare.PatchARecords(ipv4)
		}
		if newIpv6 != ipv6 {
			ipv6 = newIpv6
			cloudflare.PatchAAAARecords(ipv6)
		}
		time.Sleep(time.Second * time.Duration(cfg.Timeout))
	}
}
