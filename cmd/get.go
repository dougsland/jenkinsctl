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

// getCmd represents the show command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "Display informations about Jenkins",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("❌ requires at least one argument")
		}
		return nil
	},
}

// Connection Command
var connectionInfo = &cobra.Command{
	Use:   "connection",
	Short: "get connection info",
	RunE: func(cmd *cobra.Command, args []string) error {
		j.ServerInfo()
		return nil
	},
}

// Plugins Command
var pluginsInfo = &cobra.Command{
	Use:   "plugins",
	Short: "get all plugins active and enabled",
	RunE: func(cmd *cobra.Command, args []string) error {
		j.PluginsShow()
		return nil
	},
}

var viewsInfo = &cobra.Command{
	Use:   "views",
	Short: "get all views",
	RunE: func(cmd *cobra.Command, args []string) error {
		j.ShowViews()
		return nil
	},
}

// Build Command
var build = &cobra.Command{
	Use:   "build",
	Short: "build related commands",
}

var buildQueue = &cobra.Command{
	Use:   "queue",
	Short: "get build queue",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("⏳ Collecting build queue information...\n")
		j.ShowBuildQueue()
		return nil
	},
}

// Node Command
var nodes = &cobra.Command{
	Use:   "nodes",
	Short: "nodes related commands",
}

var nodesOffline = &cobra.Command{
	Use:   "offline",
	Short: "get nodes offline",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("⏳ Collecting node(s) information...\n")
		j.ShowNodes("offline")
		return nil
	},
}

var nodesOnline = &cobra.Command{
	Use:   "online",
	Short: "get nodes online",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("⏳ Collecting node(s) information...\n")
		j.ShowNodes("online")
		return nil
	},
}

// Var related to flags
var (
	connectionInfoBool = false
	pluginsInfoBool    = false
	viewsInfoBool      = false
	nodesBool          = false
	nodesOfflineBool   = false
	nodesOnlineBool    = false
)

func init() {
	rootCmd.AddCommand(getCmd)

	connectionInfo.Flags().BoolVarP(&connectionInfoBool,
		"connection",
		"c",
		false,
		"get connection info")

	pluginsInfo.Flags().BoolVarP(&pluginsInfoBool,
		"plugins",
		"p",
		false,
		"get all plugins actived and enabled")

	viewsInfo.Flags().BoolVarP(&viewsInfoBool,
		"views",
		"v",
		false,
		"get all views")

	nodes.Flags().BoolVarP(&nodesBool,
		"nodes",
		"n",
		false,
		"get all nodes")

	nodesOffline.Flags().BoolVarP(&nodesOfflineBool,
		"offline",
		"o",
		false,
		"get all nodes offline")

	// get
	getCmd.AddCommand(connectionInfo)
	getCmd.AddCommand(pluginsInfo)
	getCmd.AddCommand(viewsInfo)
	getCmd.AddCommand(nodes)
	getCmd.AddCommand(build)

	// nodes
	nodes.AddCommand(nodesOffline)
	nodes.AddCommand(nodesOnline)

	// build
	build.AddCommand(buildQueue)
}
