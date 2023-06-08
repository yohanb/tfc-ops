package lib

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Provider struct {
	ID         string `json:"id"`
	Type       string `json:"type"`
	Attributes struct {
		Name         string `json:"name"`
		Namespace    string `json:"namespace"`
		CreatedAt    string `json:"created-at"`
		UpdatedAt    string `json:"updated-at"`
		RegistryName string `json:"registry-name"`
	} `json:"attributes"`
	Links struct {
		Self string `json:"self"`
	} `json:"links"`
}

type ProviderList struct {
	Data []Provider `json:"data"`
}

// GetAllProviders retrieves all registry providers from Terraform Cloud and returns a list of Provider objects
func GetAllProviders(organization string) ([]Provider, error) {
	u := NewTfcUrl(fmt.Sprintf("/organizations/%s/registry-providers", organization))
	u.SetParam(paramPageSize, strconv.Itoa(pageSize))

	var allProviderData []Provider

	for page := 1; ; page++ {
		u.SetParam(paramPageNumber, strconv.Itoa(page))
		nextProviderData, err := getProviderPage(u.String())
		if err != nil {
			return []Provider{}, fmt.Errorf("error getting provider data for %s: %s", organization, err)
		}
		allProviderData = append(allProviderData, nextProviderData.Data...)

		// If there isn't a whole page of contents, then we're on the last one.
		if len(nextProviderData.Data) < pageSize {
			break
		}
	}
	return allProviderData, nil
}

func getProviderPage(url string) (ProviderList, error) {
	resp := callAPI(http.MethodGet, url, "", nil)

	defer resp.Body.Close()

	var nextProviderData ProviderList

	if err := json.NewDecoder(resp.Body).Decode(&nextProviderData); err != nil {
		return ProviderList{}, fmt.Errorf("json decode error: %s", err)
	}
	return nextProviderData, nil
}
