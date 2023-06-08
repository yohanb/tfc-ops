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

var distFolder string
var gpgKeyId string

var providerVersionPublishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish a provider version",
	Long:  `Show the values of variables with a key or value containing a certain string`,
	Args:  cobra.ExactArgs(0),
	Run: func(cmd *cobra.Command, args []string) {
		runProviderPublish()
	},
}

func init() {
	providerVersionCmd.AddCommand(providerVersionPublishCmd)
	addProviderRelatedFlags(providerVersionPublishCmd)

	providerVersionPublishCmd.Flags().StringVarP(&providerVersion, "version", "v", "",
		requiredPrefix+"provider version")
	if err := providerVersionPublishCmd.MarkFlagRequired("version"); err != nil {
		errLog.Fatalln("failed to mark 'version' as a required flag on command: " + providerVersionPublishCmd.Name())
	}

	providerVersionPublishCmd.Flags().StringVar(&providerPlatformOs, "platform-os", "",
		requiredPrefix+"provider platform-os")
	if err := providerVersionPublishCmd.MarkFlagRequired("platform-os"); err != nil {
		errLog.Fatalln("failed to mark 'platform-os' as a required flag on command: " + providerVersionPublishCmd.Name())
	}

	providerVersionPublishCmd.Flags().StringVar(&providerPlatformArch, "platform-arch", "",
		requiredPrefix+"provider platform-arch")
	if err := providerVersionPublishCmd.MarkFlagRequired("platform-arch"); err != nil {
		errLog.Fatalln("failed to mark 'platform-arch' as a required flag on command: " + providerVersionPublishCmd.Name())
	}

	providerVersionPublishCmd.Flags().StringVar(&gpgKeyId, "key-id", "",
		requiredPrefix+"gpg key id")
	if err := providerVersionPublishCmd.MarkFlagRequired("key-id"); err != nil {
		errLog.Fatalln("failed to mark 'key-id' as a required flag on command: " + providerVersionPublishCmd.Name())
	}

	providerVersionPublishCmd.Flags().StringVar(&distFolder, "dist-folder", "./dist",
		requiredPrefix+"provider dist-folder. Contains the SHA256SUM and provider archived files. Defaults to './dist'")
}

func runProviderPublish() {
	if readOnlyMode {
		fmt.Println("Read only mode enabled. No provider will be published.")
	}

	if providerVersion == "" {
		errLog.Fatal("No version specified")
	}

	if providerPlatformOs == "" {
		errLog.Fatal("No platform-os specified")
	}

	if providerPlatformArch == "" {
		errLog.Fatal("No platform-arch specified")
	}

	if gpgKeyId == "" {
		errLog.Fatal("No gpg key id specified")
	}

	found := publishProvider(
		organization,
		providerNamespace,
		providerName,
		providerVersion,
		providerPlatformOs,
		providerPlatformArch,
		gpgKeyId,
		distFolder,
	)
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

func publishProvider(org string, namespace string, name string, version string, os string, arch string, keyId string, distFolder string) bool {
	fmt.Printf(
		"Publishing provider '%s' version '%s' for platform '%s_%s' in namespace '%s'\n",
		name,
		version,
		os,
		arch,
		namespace,
	)
	if !readOnlyMode {
		err := lib.CreateProviderVersion(org, namespace, name, version, keyId)
		if err != nil {
			return false
		}
	}
	return true
}
