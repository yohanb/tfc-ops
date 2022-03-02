package lib

import (
	"net/url"
	"strconv"
)

const (
	baseURLv2 = "https://app.terraform.io/api/v2"

	pageSize = 20

	paramFilterOrganizationName = "filter[organization][name]"
	paramFilterWorkspaceID      = "filter[workspace][id]"
	paramFilterWorkspaceName    = "filter[workspace][name]"
	paramPageSize               = "page[size]"
	paramPageNumber             = "page[number]"
	paramSearchName             = "search[name]"
)

type TfcUrl struct {
	url.URL
}

func NewTfcUrl(path string) TfcUrl {
	newURL, _ := url.Parse(baseURLv2 + path)
	v := url.Values{}
	newURL.RawQuery = v.Encode()
	tfcUrl := TfcUrl{
		URL: *newURL,
	}
	tfcUrl.SetParam(paramPageSize, strconv.Itoa(pageSize))
	return tfcUrl
}

func (t *TfcUrl) SetParam(name, value string) {
	values := t.Query()
	values.Set(name, value)
	t.RawQuery = values.Encode()
}
