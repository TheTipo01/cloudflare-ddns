package main

import (
	"github.com/kkyr/fig"
	"github.com/nylone/cloudflare-ddns/src/cloudflare"
	"plugin"
	"sync"
)

const (
	configFile = "config/config.yml"
)

type config struct {
	Backend string `fig:"backend" validate:"required"`
}

func main() {
	// load config file
	var cfg config
	err := fig.Load(&cfg, fig.File(configFile))
	if err != nil {
		panic(err)
	}

	// find correct plugin
	var pluginPath string
	switch cfg.Backend {
	case "ipify":
		pluginPath = "plugins/ipify/ipify.so"
	case "skywifi":
		pluginPath = "plugins/skywifi/skywifi.so"
	case "openwrt":
		pluginPath = "plugins/openwrt/openwrt.so"
	case "vodafone":
		pluginPath = "plugins/vodafone/vodafone.so"
	default:
		panic("No valid plugin selected.")
	}

	// load plugin
	p, err := plugin.Open(pluginPath)
	if err != nil {
		panic(err)
	}
	f, err := p.Lookup("GetIPs")
	if err != nil {
		panic(err)
	}
	getIPs := f.(func(bool, bool) (string, string, error))

	// load cloudflare mappings
	err = cloudflare.LoadMappings()
	if err != nil {
		panic(err)
	}

	ipv4, ipv6, err := getIPs(true, true)
	if err != nil {
		wg := sync.WaitGroup{}
		defer wg.Wait()
		cloudflare.PatchARecords(ipv4, &wg)
		cloudflare.PatchAAAARecords(ipv6, &wg)
	}
}
