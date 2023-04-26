package main

import (
	"encoding/json"
	"github.com/kkyr/fig"
	"io"
	"net/http"
	"strings"
)

const (
	configFile = "config/vodafone.yml"
)

var (
	replacer *strings.Replacer
	req      *http.Request
	cfg      config
)

type config struct {
	Router string `fig:"router" validate:"required"`
}

type StationJSON struct {
	WanIP4Addr string `json:"wan_ip4_addr"`
	WanIP6Addr string `json:"wan_ip6_addr"`
}

func init() {
	err := fig.Load(&cfg, fig.File(configFile))
	if err != nil {
		panic(err)
	}

	replacer = strings.NewReplacer("[", "{", "]", "}", "{", "", "}", "")
	req, _ = http.NewRequest("GET", "http://"+cfg.Router+"/data/user_lang.json", nil)
	// Add Accept-Language header, otherwise the modem will throw bad requests at us
	req.Header.Set("Accept-Language", "it-IT")
}

// GetIPs returns public IP address of the machine from the specified endpoint
func GetIPs(loadV4, loadV6 bool) (string, string, error) {
	var (
		out        StationJSON
		ipv4, ipv6 string
	)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}

	b, _ := io.ReadAll(resp.Body)
	_ = resp.Body.Close()

	// The JSON is given to us in an array. We parse that and remove the brackets, and add them only at the end
	err = json.Unmarshal([]byte(replacer.Replace(string(b))), &out)
	if err != nil {
		return "", "", err
	}

	if loadV4 && out.WanIP4Addr != "N/A" {
		ipv4 = out.WanIP4Addr
	}
	if loadV6 && out.WanIP6Addr != "N/A" {
		ipv6 = out.WanIP6Addr
	}

	return ipv4, ipv6, nil
}
