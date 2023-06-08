package lib

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// callAPI creates a http.Request object, attaches headers to it and makes the
// requested api call.
func callAPI(method, url, postData string, headers map[string]string) *http.Response {
	var err error
	var req *http.Request

	proxy := os.Getenv("HTTP_PROXY")
	if proxy != "" {
		fmt.Println("Using proxy: " + proxy)
	}

	if postData != "" {
		req, err = http.NewRequest(method, url, strings.NewReader(postData))
	} else {
		req, err = http.NewRequest(method, url, nil)
	}

	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	req.Header.Set("Authorization", "Bearer "+config.token)
	req.Header.Set("Content-Type", "application/vnd.api+json")

	for key, val := range headers {
		req.Header.Set(key, val)
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else if resp.StatusCode >= 300 {
		bodyBytes, _ := ioutil.ReadAll(resp.Body)
		fmt.Println(fmt.Sprintf(
			"API returned an error.\n\tMethod: %s\n\tURL: %s\n\tCode: %v\n\tStatus: %s\n\tRequest Body: %s\n\tResponse Body: %s",
			method, url, resp.StatusCode, resp.Status, postData, bodyBytes))
		os.Exit(1)
	}

	return resp
}
