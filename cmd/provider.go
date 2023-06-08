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
	"github.com/spf13/cobra"
)

var (
	providerName         string
	providerNamespace    string
	providerVersion      string
	providerPlatformOs   string
	providerPlatformArch string
)

func addProviderRelatedFlags(command *cobra.Command) {
	command.Flags().StringVarP(&providerName, "name", "n", "",
		requiredPrefix+"provider name")
	if err := command.MarkFlagRequired("name"); err != nil {
		errLog.Fatalln("failed to mark 'name' as a required flag on command: " + command.Name())
	}

	command.Flags().StringVarP(&providerNamespace, "namespace", "s", "",
		requiredPrefix+"provider namespace")
	if err := command.MarkFlagRequired("namespace"); err != nil {
		errLog.Fatalln("failed to mark 'namespace' as a required flag on command: " + command.Name())
	}
}

// providerCmd represents the top level command for variables
var providerCmd = &cobra.Command{
	Use:   "provider",
	Short: "Configure registry provider",
	Long:  `Top level command to configure your private registry provider in Terraform Cloud`,
	Args:  cobra.MinimumNArgs(1),
}

func init() {
	rootCmd.AddCommand(providerCmd)
	addGlobalFlags(providerCmd)
}
