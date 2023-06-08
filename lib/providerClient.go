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

type ProviderVersionAttributes struct {
	Version   string   `json:"version"`
	KeyId     string   `json:"key-id"`
	Protocols []string `json:"protocols"`
}

type ProviderVersion struct {
	Type       string                    `json:"type"`
	Attributes ProviderVersionAttributes `json:"attributes"`
}

type ProviderVersionList struct {
	Data []ProviderVersion `json:"data"`
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

func GetAllProviderVersions(organization string, providerNamespace string, providerName string) ([]ProviderVersion, error) {
	u := NewTfcUrl(fmt.Sprintf("/organizations/%s/registry-providers/private/%s/%s/versions", organization, providerNamespace, providerName))
	u.SetParam(paramPageSize, strconv.Itoa(pageSize))

	var allProviderVersionData []ProviderVersion

	for page := 1; ; page++ {
		u.SetParam(paramPageNumber, strconv.Itoa(page))
		nextProviderVersionData, err := getProviderVersionPage(u.String())
		if err != nil {
			return []ProviderVersion{}, fmt.Errorf("error getting provider version data for org %s, providerNamespace %s, name %s: %s", organization, providerNamespace, providerName, err)
		}
		allProviderVersionData = append(allProviderVersionData, nextProviderVersionData.Data...)

		// If there isn't a whole page of contents, then we're on the last one.
		if len(nextProviderVersionData.Data) < pageSize {
			break
		}
	}
	return allProviderVersionData, nil
}

func CreateProviderVersion(organization string, providerNamespace string, providerName string, providerVersion string, gpgKeyId string) error {
	u := NewTfcUrl(fmt.Sprintf(
		"/organizations/%s/registry-providers/private/%s/%s/versions",
		organization,
		providerNamespace,
		providerName,
	))

	version := ProviderVersion{
		Type: "registry-provider-versions",
		Attributes: ProviderVersionAttributes{
			Version:   providerVersion,
			KeyId:     gpgKeyId,
			Protocols: []string{"5.0"}},
	}

	data := struct {
		Data ProviderVersion `json:"data"`
	}{version}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	resp := callAPI(http.MethodPost, u.String(), string(jsonData), nil)
	_ = resp.Body.Close()

	return nil
}

func DeleteProviderVersion(organization string, providerNamespace string, providerName string, providerVersion string) {
	u := NewTfcUrl(fmt.Sprintf(
		"/organizations/%s/registry-providers/private/%s/%s/versions/%s",
		organization,
		providerNamespace,
		providerName,
		providerVersion))

	resp := callAPI(http.MethodDelete, u.String(), "", nil)
	_ = resp.Body.Close()
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

func getProviderVersionPage(url string) (ProviderVersionList, error) {
	resp := callAPI(http.MethodGet, url, "", nil)

	defer resp.Body.Close()

	var nextProviderVersionData ProviderVersionList

	if err := json.NewDecoder(resp.Body).Decode(&nextProviderVersionData); err != nil {
		return ProviderVersionList{}, fmt.Errorf("json decode error: %s", err)
	}
	return nextProviderVersionData, nil
}
