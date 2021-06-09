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
	"github.com/dougsland/jenkinscli/jenkins"
	"github.com/spf13/cobra"
)

var j jenkins.Jenkins

// connectionCmd represents the connection command
var connectionCmd = &cobra.Command{
	Use:   "connection [OPTIONS]",
	Short: "Options related to connections",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("requires at least one arg")
		}
		return nil
	},
}

var showConnCmd = &cobra.Command{
	Use:   "show",
	Short: "show connection info",
	RunE: func(cmd *cobra.Command, args []string) error {
		err := j.Init()
		if err != nil {
			return err
		}
		j.ServerInfo()
		return nil
	},
}

func init() {
	var showConnCmdBool = false
	j = jenkins.Jenkins{}

	rootCmd.AddCommand(connectionCmd)
	showConnCmd.Flags().BoolVarP(&showConnCmdBool,
		"show",
		"s",
		false,
		"show connection info")
	connectionCmd.AddCommand(showConnCmd)
}
