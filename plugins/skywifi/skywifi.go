package main

import (
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func main() {
	page, _ := getIP()
	print(page)
}

// return public IP address of the machine from the specified endpoint
func getIP() (string, error) {
	// http client able to work with cookies
	jar, _ := cookiejar.New(nil)
	client := &http.Client{
		Jar: jar,
	}

	// authenticate to webserver
	{
		v := url.Values{}
		v.Set("username", "admin")
		v.Add("password", "password")
		request, err := http.NewRequest("POST", "http://192.168.0.1/check.jst", strings.NewReader(v.Encode()))
		if err != nil {
			return "", err
		}
		response, err := client.Do(request)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()
	}

	// request connection info page
	{
		request, err := http.NewRequest("GET", "http://192.168.0.1/network_setup.jst", nil)
		if err != nil {
			return "", err
		}
		response, err := client.Do(request)
		if err != nil {
			return "", err
		}
		defer response.Body.Close()
		doc, err := goquery.NewDocumentFromReader(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		return string(doc.Find(".value").Text()), nil
	}

}
