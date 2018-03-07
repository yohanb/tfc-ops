package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type TFMeta struct {
	Total int `json:"total"`
}

type TFEnv struct {
	Username string `json:"username"`
	Name     string `json:"name"`
}

type TFState struct {
	UpdatedAt   string `json:"updated_at"`
	Environment TFEnv  `json:"environment"`
}

type TFAllStates struct {
	States []TFState `json:"states"`
	Meta   TFMeta    `json:"meta"`
}

func getJsonFromFile(jsonFile string) TFAllStates {
	raw, err := ioutil.ReadFile(jsonFile)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var contents TFAllStates
	json.Unmarshal(raw, &contents)
	return contents
}

/*
 * @param jsonFile The path and/or file name of a json file the contents
 *        of which match what is returned from the terraform api
 * @return a slice of strings -the environment names
 */
func GetAllEnvNamesFromJson(jsonFile string) []string {
	envNames := []string{}
	allStates := getJsonFromFile(jsonFile)

	for _, nextState := range allStates.States {
		envNames = append(envNames, nextState.Environment.Name)
	}

	return envNames
}

/*
 * @param tfToken The user's Terraform Enterprise Token
 * @return a slice of strings - the environment names from the v1 api
 */
func GetAllEnvNamesFromV1API(tfToken string) []string {
	baseURL := "https://atlas.hashicorp.com/api/v1/terraform/state?page="
	names := []string{}

	for page := 1; ; page++ {
		url := fmt.Sprintf("%s%d", baseURL, page)

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		req.Header.Set("X-Atlas-Token", tfToken)

		client := &http.Client{}

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		defer resp.Body.Close()

		var statesFromPage TFAllStates

		if err := json.NewDecoder(resp.Body).Decode(&statesFromPage); err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}

		for _, nextState := range statesFromPage.States {
			names = append(names, nextState.Environment.Name)
		}

		// If there isn't a whole page of contents, then we're on the last one.
		if len(statesFromPage.States) < 20 {
			break
		}
	}

	return names
}
