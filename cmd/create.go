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
	"os"
)

// createCmd represents the create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a resource in Jenkins",
}

var createNode = &cobra.Command{
	Use:   "node",
	Short: "create a node",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("In the TODO list")
		return nil
	},
}

func detectViewType(view string) string {
	viewSelected := ""
	switch view {
	case "LIST_VIEW":
		viewSelected = "hudson.model.ListView"
		break
	case "NESTED_VIEW":
		viewSelected = "hudson.plugins.nested_view.NestedView"
		break
	case "MY_VIEW":
		viewSelected = "hudson.model.MyView"
		break
	case "DASHBOARD_VIEW":
		viewSelected = "hudson.plugins.view.dashboard.Dashboard"
		break
	case "PIPELINE_VIEW":
		viewSelected = "au.com.centrumsystems.hudson.plugin.buildpipeline.BuildPipelineView"
		break
	default:
		fmt.Println("error: use only views supported: LIST_VIEW, NESTED_VIEW, MY_VIEW, DASHBOARD_VIEW, PIPELINE_VIEW")
		os.Exit(1)
	}

	return viewSelected
}

var createView = &cobra.Command{
	Use:   "view",
	Short: "create a view",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("❌ requires at least two arguments: viewName viewType (LIST_VIEW, NESTED_VIEW, MY_VIEW, DASHBOARD_VIEW, PIPELINE_VIEW")
		}

		fmt.Printf("⏳ Creating view %s...\n", args[0])
		err := jenkinsMod.CreateView(args[0], detectViewType(args[1]))
		if err != nil {
			fmt.Printf("unable to create the view: %s - err: %s \n", args[1], err)
			os.Exit(1)
		}
		return nil
	},
}

var createJob = &cobra.Command{
	Use:   "job",
	Short: "create a job",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 2 {
			return errors.New("❌ requires at least two arguments: xmlFile JobName")
		}

		fmt.Printf("⏳ Creating the job %s...\n", args[0])
		err := jenkinsMod.CreateJob(args[0], args[1])
		if err != nil {
			fmt.Printf("unable to create the job: %s - err: %s \n", args[1], err)
			os.Exit(1)
		}
		fmt.Printf("Created job: %s\n", args[1])
		return nil
	},
}

func init() {
	rootCmd.AddCommand(createCmd)
	createCmd.AddCommand(createJob)
	createCmd.AddCommand(createNode)
	createCmd.AddCommand(createView)
}
