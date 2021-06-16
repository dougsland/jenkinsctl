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
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "download related commands",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("download called")
	},
}

var artifactsCmd = &cobra.Command{
	Use:   "artifacts",
	Short: "download artifacts",
	Run: func(cmd *cobra.Command, args []string) {

		if len(args) != 3 {
			fmt.Println("❌ requires at three arguments [JOB_NAME BUILD_ID PATH_TO_SAVE_ARTIFACTS]")
			os.Exit(1)
		}

		err := j.DownloadArtifacts(args[0], args[1], args[2])
		if err != nil {
			fmt.Println("cannot download artifacts: %s", err)
			os.Exit(1)
		}
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)
	downloadCmd.AddCommand(artifactsCmd)

}
