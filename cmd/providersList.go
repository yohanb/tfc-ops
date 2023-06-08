// Copyright © 2018-2022 SIL International
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/silinternational/tfc-ops/lib"
	"github.com/spf13/cobra"
)

var providersListCmd = &cobra.Command{
	Use:   "list",
	Short: "Report on providers",
	Long:  `Show the values of variables with a key or value containing a certain string`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Getting list of providers ...")
		allData, err := lib.GetAllProviders(organization)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		for _, provider := range allData {
			fmt.Println(provider.Attributes.Name)
		}
	},
}

func init() {
	providersCmd.AddCommand(providersListCmd)
}
