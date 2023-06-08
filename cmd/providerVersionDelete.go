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

var providerVersionDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a provider version",
	Long:  `Deletes the specified provider version`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runVersionDelete()
	},
}

func init() {
	providerVersionCmd.AddCommand(providerVersionDeleteCmd)
	addProviderRelatedFlags(providerVersionDeleteCmd)

	providerVersionDeleteCmd.Flags().StringVarP(&providerVersion, "version", "v", "",
		requiredPrefix+"provider version")
	if err := providerVersionDeleteCmd.MarkFlagRequired("version"); err != nil {
		errLog.Fatalln("failed to mark 'version' as a required flag on command: " + providerVersionDeleteCmd.Name())
	}
}

func runVersionDelete() {
	if readOnlyMode {
		fmt.Println("Read only mode enabled. No provider version will be deleted.")
	}

	if providerVersion == "" {
		errLog.Fatal("No version specified")
	}

	found := deleteProviderVersion(organization, providerNamespace, providerName, providerVersion)
	if !found {
		errLog.Fatalf(
			"Version '%s' not found for provider '%s' in namespace '%s'\n",
			version,
			providerName,
			providerNamespace,
		)
	}
	return
}

func deleteProviderVersion(org string, namespace string, name string, version string) bool {
	fmt.Printf(
		"Deleting version '%s' from provider '%s' in namespace '%s'\n",
		version,
		name,
		namespace,
	)
	if !readOnlyMode {
		lib.DeleteProviderVersion(org, namespace, name, version)
	}
	return true
}
