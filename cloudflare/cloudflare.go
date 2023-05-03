package cloudflare

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/kkyr/fig"
)

const (
	baseAPIUrl = "https://api.cloudflare.com/client/v4/zones/"
	configFile = "config/cloudflare.yml"
)

var (
	cfg   config                                          // config file mappings
	z2v4r map[string][]string = make(map[string][]string) // maps a zone to an array of its records
	z2v6r map[string][]string = make(map[string][]string) // maps a zone to an array of its records
)

type config struct {
	Token string `fig:"token" validate:"required"`
	Zones map[string]struct {
		V4Records map[string]interface{} `fig:"v4-records" validate:"required"`
		V6Records map[string]interface{} `fig:"v6-records" validate:"required"`
	} `fig:"zones" validate:"required"`
}

// Loads the config file
func init() {
	err := fig.Load(&cfg, fig.File(configFile))
	if err != nil {
		panic(err)
	}
}

// Reloads information from cloudflare about zones and dns records available
func LoadMappings() error {
	zones, err := getZones()
	if err != nil {
		return err
	}
	for _, z := range zones {
		cfgRecords, ok := cfg.Zones[z.Name]
		if ok {
			records, err := getRecords(z)
			if err != nil {
				return err
			}
			for _, r := range records {
				switch r.Type {
				case "A":
					_, ok := cfgRecords.V4Records[r.Name]
					if ok {
						z2v4r[z.ID] = append(z2v4r[z.ID], r.ID)
					}
				case "AAAA":
					_, ok := cfgRecords.V6Records[r.Name]
					if ok {
						z2v6r[z.ID] = append(z2v6r[z.ID], r.ID)
					}
				}
			}
		}
	}
	return nil
}

// Update A records on cloudflare.
func PatchARecords(ipv4 string) []error {
	mutexError := sync.Mutex{}
	var errors []error
	for z, r := range z2v4r {
		for _, r := range r {
			zone, record := z, r
			go func() {
				err := patchRecord(zone, record, ipv4)
				if err != nil {
					mutexError.Lock()
					defer mutexError.Unlock()
					errors = append(errors, err)
				}
			}()
		}
	}
	return errors
}

// Update AAAA records on cloudflare.
func PatchAAAARecords(ipv6 string) []error {
	mutexError := sync.Mutex{}
	var errors []error
	for z, r := range z2v4r {
		for _, r := range r {
			zone, record := z, r
			go func() {
				err := patchRecord(zone, record, ipv6)
				if err != nil {
					mutexError.Lock()
					defer mutexError.Unlock()
					errors = append(errors, err)
				}
			}()
		}
	}
	return errors
}

type zone struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type record struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Type string `json:"type"`
}

// Retrieves all zones accessible from the configured token.
//
// Requires a valid token to have been setup.
func getZones() ([]zone, error) {
	type apiResponse struct {
		Result  []zone `json:"result"`
		Success bool   `json:"success"`
	}

	request, err := http.NewRequest("GET", baseAPIUrl, nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("authorization", "Bearer "+cfg.Token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var parsedResponse apiResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	if err != nil {
		return nil, err
	}
	if parsedResponse.Success {
		return parsedResponse.Result, nil
	} else {
		return nil, errors.New("could not load zones from cloudflare's api")
	}
}

// Retrieves all records for a given zone.
//
// Must be given the zone to fetch for.
func getRecords(z zone) ([]record, error) {
	type apiResponse struct {
		Result  []record `json:"result"`
		Success bool     `json:"success"`
		Errors  []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}

	request, err := http.NewRequest("GET", strings.Join([]string{baseAPIUrl, z.ID, "/dns_records/"}, ""), nil)
	if err != nil {
		return nil, err
	}
	request.Header.Add("authorization", "Bearer "+cfg.Token)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	var parsedResponse apiResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	if err != nil {
		return nil, err
	}
	if parsedResponse.Success {
		return parsedResponse.Result, nil
	} else {
		return nil, errors.New(fmt.Sprint(parsedResponse.Errors))
	}
}

// Tells cloudflare what the new ip of a record is.
//
// Must be given the zone id {zid}, the record id {rid} and the new ip {ip}.
func patchRecord(zid string, rid string, ip string) error {
	type apiResponse struct {
		Success bool `json:"success"`
		Errors  []struct {
			Code    int    `json:"code"`
			Message string `json:"message"`
		} `json:"errors"`
	}

	request, err := http.NewRequest("PATCH", strings.Join([]string{baseAPIUrl, zid, "/dns_records/", rid}, ""),
		strings.NewReader("{\"content\":\""+ip+"\"}"))
	if err != nil {
		return err
	}
	request.Header.Add("authorization", "Bearer "+cfg.Token)
	request.Header.Add("content-type", "application/json")
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()
	var parsedResponse apiResponse
	err = json.NewDecoder(response.Body).Decode(&parsedResponse)
	if err != nil {
		return err
	}
	if !parsedResponse.Success {
		return errors.New(fmt.Sprint(parsedResponse.Errors))
	}
	return nil
}
