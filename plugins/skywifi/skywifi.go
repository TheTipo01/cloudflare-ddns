package main

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"regexp"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/kkyr/fig"
)

const (
	configFile = "config/skywifi.yml"
)

var (
	v4Regex = regexp.MustCompile("([\\d.]+)")
	v6Regex = regexp.MustCompile("([a-f0-9:]+:+)+[a-f0-9]+")

	// http client able to work with cookies
	jar, _ = cookiejar.New(nil)
	client = &http.Client{
		Jar: jar,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}
	cfg config // config file mappings

)

type config struct {
	Username string `fig:"username" validate:"required"`
	Password string `fig:"password" validate:"required"`
	Router   string `fig:"router" validate:"required"`
}

// Loads the config file
func init() {
	err := fig.Load(&cfg, fig.File(configFile))
	if err != nil {
		panic(err)
	}
}

// GetIPs returns public IP address of the machine from the specified endpoint
func GetIPs(loadV4, loadV6 bool) (string, string, error) {
	// authenticate to webserver
	v := url.Values{}
	v.Set("username", cfg.Username)
	v.Add("password", cfg.Password)
	request, err := http.NewRequest("POST", fmt.Sprint("http://", cfg.Router, "/check.jst"), strings.NewReader(v.Encode()))
	if err != nil {
		return "", "", err
	}
	response, err := client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	// request connection info page
	var doc *goquery.Document
	request, err = http.NewRequest("GET", fmt.Sprint("http://", cfg.Router, "/network_setup.jst"), nil)
	if err != nil {
		return "", "", err
	}
	response, err = client.Do(request)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()
	doc, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	// scrape and find ipv4 and ipv6 addresses from the document
	var ipv4, ipv6 = "", ""
	if loadV4 {
		ipv4 = doc.Find("#wanip4").Next().Text()
		ipv4 = v4Regex.FindAllString(ipv4, 1)[0]
	}
	if loadV6 {
		routerIpv6 := doc.Find("#wanip6").Next().Text()
		routerIpv6 = v6Regex.FindAllString(routerIpv6, 1)[0]
		ipv6, err = findPublicIPv6(routerIpv6)
		if err != nil {
			return "", "", err
		}
	}
	return ipv4, ipv6, nil
}

func findPublicIPv6(routerv6 string) (string, error) {
	// extract current netmask from the router ip
	_, network, err := net.ParseCIDR(routerv6 + "/64")
	if err != nil {
		return "", err
	}

	// get list of available addresses
	addr, err := net.InterfaceAddrs()
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	for _, addr := range addr {
		if ip, ok := addr.(*net.IPNet); ok && network.Contains(ip.IP) {
			return ip.IP.String(), nil
		}
	}

	return "", errors.New("no public ipv6 was found for router address: " + routerv6)
}
