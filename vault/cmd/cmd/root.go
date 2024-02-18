// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
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

package vault

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

// RootVaultCmd is the root command for vault cli m,,l
var RootVaultCmd = &cobra.Command{
	Use:   "vault",
	Short: "Vault helps with managing your secret configuration on the go.",
	Long:  `This application helps you manage your secrets while encrypting them using AES encryption.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {

		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := RootVaultCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	//cobra.OnInitialize(initConfig)
}