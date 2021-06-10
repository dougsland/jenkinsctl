/*
Copyright © 2021 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"errors"
	"fmt"

	"github.com/spf13/cobra"
)

// setCmd represents the set command
var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set a resource regarding jenkinsctl",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("❌ requires at least one argument")
		}
		return nil
	},
}

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "set config path",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("....")
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.AddCommand(configCmd)
}
