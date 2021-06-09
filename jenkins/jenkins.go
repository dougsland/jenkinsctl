package jenkins

import (
	"context"
	"errors"
	"fmt"
	"github.com/bndr/gojenkins"
	"github.com/spf13/viper"
	"net/http"
	"os"
)

// Jenkins connection object
type Jenkins struct {
	Instance *gojenkins.Jenkins
	Server   string
	Username string
	Token    string
	Context  context.Context
}

// Config is focused in the configuration json file
type Config struct {
	Server         string `mapstructure: Server`
	Username       string `mapstructure: Username`
	Admuser        string `mapstructure: Admuser`
	Token          string `mapstructure: Token`
	ConfigPath     string
	ConfigFileName string
	ConfigFullPath string
}

// SetConfigPath set the default config path
//
// Args:
//
// Returns
//	string or error
func (j *Config) SetConfigPath() error {
	home := os.Getenv("HOME")
	if len(home) == 0 {
		return errors.New("cannot get $HOME env var")
	}
	j.ConfigPath = home + "/.config/" + "jenkinscli/"
	j.ConfigFileName = "config.json"
	j.ConfigFullPath = j.ConfigPath + j.ConfigFileName
	return nil
}

// CheckIfExists check if file exists
//
// Args:
//	path - string
//
// Returns
//	error
func (j *Config) CheckIfExists() error {
	var err error
	if _, err = os.Stat(j.ConfigFullPath); err == nil {
		return nil

	}
	return err
}

// LoadConfig read the JSON configuration from specified file
//
// Example file:
//
// $HOME/.config/jenkinscli/config.json
//
//
// Args:
//
// Returns
//	nil or error
func (j *Config) LoadConfig() (config Config, err error) {

	viper.AddConfigPath(j.ConfigPath)
	viper.SetConfigName(j.ConfigFileName)
	viper.SetConfigType("json")
	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(&config)
	return
}

// PluginsShow show all plugins installed and enabled
//
// Returns
//	nil or error
func (j *Jenkins) PluginsShow() {
	p, _ := j.Instance.GetPlugins(j.Context, 1)

	if len(p.Raw.Plugins) > 0 {
		fmt.Printf("Plugins Activated and Enabled 🚀\n")
		for _, p := range p.Raw.Plugins {
			if len(p.LongName) > 0 && p.Active && p.Enabled {
				fmt.Printf("    %s - %s ✅\n", p.LongName, p.Version)
			}
		}
	}
}

// ShowBuildQueue show the Build Queue
//
// Args:
//
// Returns
//
// TIP: Meaning of collors:
// https://github.com/jenkinsci/jenkins/blob/5e9b451a11926e5b42d4a94612ca566de058f494/core/src/main/java/hudson/model/BallColor.java#L56
func (j *Jenkins) ShowBuildQueue() {
	fmt.Printf("⏳ Collecting build queue information...\n\n")

	queue, _ := j.Instance.GetQueue(j.Context)
	totalTasks := 0
	for i, item := range queue.Raw.Items {
		fmt.Printf("URL: %s\n", item.Task.URL)
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Name: %s\n", item.Task.Name)
		fmt.Printf("Pending: %v\n", item.Pending)
		fmt.Printf("Stuck: %v\n", item.Stuck)

		switch item.Task.Color {
		case "red":
			fmt.Printf("Status: ❌ Failed\n")
			break
		case "red_anime":
			fmt.Printf("Status: ⏳ In Progress\n")
			break
		case "notbuilt":
			fmt.Printf("Status: 🚧 Not Build\n")
			break
		}
		fmt.Printf("Why: %s\n", item.Why)
		fmt.Printf("----------------")
		fmt.Printf("\n")
		totalTasks = i + 1
	}
	fmt.Printf("\nNumber of tasks in the build queue: %d\n", totalTasks)
}

func (j *Jenkins) ShowNodesOnline() error {
	return j.ShowNodes("online")
}

func (j *Jenkins) ShowNodesOffline() error {
	return j.ShowNodes("offline")
}

// ShowViews
func (j *Jenkins) ShowViews(showStatus string) error {
	views, err := j.Instance.GetAllViews(j.Context)
	if err != nil {
		return err
	}
	for _, view := range views {
		fmt.Printf("%s\n", view.Raw.Name)
		fmt.Printf("%s\n", view.Raw.URL)
	}
	return nil
}

// ShowNodes show all plugins installed and enabled
//
// Args:
//	showStatus - show only the
//
// Returns
//	nil or error
func (j *Jenkins) ShowNodes(showStatus string) error {
	nodes, err := j.Instance.GetAllNodes(j.Context)
	if err != nil {
		return err
	}
	for _, node := range nodes {
		// Fetch Node Data
		node.Poll(j.Context)

		switch showStatus {

		case "offline":
			if node.Raw.Offline || node.Raw.TemporarilyOffline {
				fmt.Printf("❌ %s - offline\n", node.GetName())
				fmt.Printf("Reason: %s\n\n", node.Raw.OfflineCauseReason)
			}

		case "online":
			if !node.Raw.Offline {
				fmt.Printf("✅ %s - online\n", node.GetName())
			}
			if node.Raw.Idle {
				fmt.Printf("😴 %s - idle\n", node.GetName())
			}
		}
	}
	return nil
}

// Init will initilialize connection with jenkins server
//
// Args:
//
// Returns
//
func (j *Jenkins) Init() {
	// Init config file
	jenkinsConfig := Config{}
	jenkinsConfig.SetConfigPath()
	config, err := jenkinsConfig.LoadConfig()
	if err != nil {
		fmt.Printf("cannot load config: %s\n", err)
		os.Exit(1)
	}

	fmt.Println("===========")
	fmt.Println(config.Admuser)

	j.Instance = gojenkins.CreateJenkins(
		nil,
		config.Server,
		config.Admuser,
		config.Token)
}

// ServerInfo will show information regarding the server
//
// Args:
//
func (j *Jenkins) ServerInfo() {
	j.Instance.Info(j.Context)
	fmt.Printf("✅ Connected with: %s\n", j.Username)
	fmt.Printf("✅ Server: %s\n", j.Server)
	fmt.Printf("✅ Version: %s\n", j.Instance.Version)
}

// serverReachable will do validation if the jenkins server
// is reachable
//
// Args:
//	string - Jenkins url
//
// Returns
//	nil or error
func serverReachable(url string) error {
	_, err := http.Get(url)
	if err != nil {
		return err
	}
	fmt.Println("Server reachable...")
	return nil

}

/*
func main() {

	// Jenkins Connection object
	jenkins := Jenkins{
		nil,
		jenkinsConfig.Server,
		jenkinsConfig.Username,
		jenkinsConfig.Token,
		context.Background(),
	}

	jenkins.Init()
	// Check if the Jenkins Server is reachable
	err = serverReachable(jenkins.Server)
	if err != nil {
		fmt.Printf("%s\n", err)
		os.Exit(1)
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	jenkins.ServerInfo()
		if os.Args[1] == "build" {
			buildCmd.Parse(os.Args[2:])
			if *buildQueueShow {
				jenkins.ShowBuildQueue()
			} else {
				fmt.Printf("❌ unknown build flag\n")
				usage()
			}
		} else if os.Args[1] == "connection" {
			connectionCmd.Parse(os.Args[2:])
			if *connectionConn {
			} else {
				fmt.Printf("❌ unknown connection flag\n")
				usage()
			}
		} else if os.Args[1] == "plugins" {
			pluginsCmd.Parse(os.Args[2:])
			if *pluginsShow {
				jenkins.PluginsShow()
			} else {
				fmt.Printf("❌ unknown plugins flag\n")
				usage()
			}
		} else if os.Args[1] == "nodes" {
			nodesCmd.Parse(os.Args[2:])

			fmt.Printf("⏳ Collecting node(s) information...\n")
			if *nodesOnline {
				jenkins.ShowNodesOnline()
			} else if *nodesOffline {
				jenkins.ShowNodesOffline()
			} else {
				fmt.Printf("❌ unknown node flag\n")
				usage()
				os.Exit(1)
			}
		} else {
			fmt.Printf("❌ unknown flag\n")
			usage()
			os.Exit(1)
		}
}
func usage() {
	fmt.Printf("usage: %s [OPTION]...\n", os.Args[0])

	fmt.Printf("OPTIONS:\n")
	fmt.Printf("\tconnection\t\tConnection Information\n")
	fmt.Printf("\t\t --show\t\tShow connection info\n")

	fmt.Printf("\n\tbuild	\t\tBuild Information\n")
	fmt.Printf("\t\t --queue\tShow the current build queue\n")

	fmt.Printf("\n\tnodes	\t\tJenkins Nodes\n")
	fmt.Printf("\t\t --online\tShow online nodes\n")
	fmt.Printf("\t\t --offline\tShow offline nodes\n")

	fmt.Printf("\n\tplugins	\t\tPlugins Information\n")
	fmt.Printf("\t\t --show\t\tShow plugins activated and enabled\n")

	fmt.Printf("\n\nExamples:\n")

	fmt.Printf("\t%s connection --show\n", os.Args[0])
	fmt.Printf("\n")

	fmt.Printf("\t%s build --queue\n", os.Args[0])
	fmt.Printf("\n")

	fmt.Printf("\t%s nodes --online\n", os.Args[0])
	fmt.Printf("\t%s nodes --offline\n", os.Args[0])

	fmt.Printf("\n")
	fmt.Printf("\t%s plugins --show\n", os.Args[0])
}
*/
