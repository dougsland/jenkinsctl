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
	Instance    *gojenkins.Jenkins
	Server      string
	JenkinsUser string
	Token       string
	Context     context.Context
}

// Config is focused in the configuration json file
type Config struct {
	Server      string `mapstructure: Server`
	JenkinsUser string `mapstructure: JenkinsUser`
	//	Admuser        string `mapstructure: Admuser`
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
		return errors.New("❌ cannot get $HOME env var")
	}
	j.ConfigPath = home + "/.config/" + "jenkinsctl/"
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
// $HOME/.config/jenkinsctl/config.json
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
func (j *Jenkins) ShowBuildQueue() error {
	err := serverReachable(j.Server)
	if err != nil {
		return errors.New("❌ jenkins server unreachable: " + j.Server)
	}

	queue, _ := j.Instance.GetQueue(j.Context)
	totalTasks := 0
	for i, item := range queue.Raw.Items {
		fmt.Printf("URL: %s\n", item.Task.URL)
		fmt.Printf("ID: %d\n", item.ID)
		fmt.Printf("Name: %s\n", item.Task.Name)
		fmt.Printf("Pending: %v\n", item.Pending)
		fmt.Printf("Stuck: %v\n", item.Stuck)

		j.ShowStatus(item.Task.Color)
		fmt.Printf("Why: %s\n", item.Why)
		fmt.Printf("----------------")
		fmt.Printf("\n")
		totalTasks = i + 1
	}
	fmt.Printf("\nNumber of tasks in the build queue: %d\n", totalTasks)

	return nil
}

// ShowStatus
func (j *Jenkins) ShowStatus(object string) {

	switch object {
	case "red":
		fmt.Printf("Status: ❌ Failed\n")
		break
	case "red_anime":
		fmt.Printf("Status: ⏳ In Progress\n")
		break
	case "notbuilt":
		fmt.Printf("Status: 🚧 Not Build\n")
		break
	default:
		fmt.Printf("Status: %s\n", object)
	}
}

// ShowViews
func (j *Jenkins) ShowViews() error {
	err := serverReachable(j.Server)
	if err != nil {
		return errors.New("❌ jenkins server unreachable: " + j.Server)
	}

	views, err := j.Instance.GetView(j.Context, "All")
	if err != nil {
		fmt.Println("erro")
		fmt.Println(err)
		return err
	}
	fmt.Println(views)
	for _, view := range views.Raw.Jobs {
		fmt.Printf("%s\n", view.Name)
		fmt.Printf("%s\n", view.Url)
		j.ShowStatus(view.Color)
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
	err := serverReachable(j.Server)
	if err != nil {
		return errors.New("❌ jenkins server unreachable: " + j.Server)
	}

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
func (j *Jenkins) Init() error {
	// Init config file
	jenkinsConfig := Config{}
	jenkinsConfig.SetConfigPath()
	config, err := jenkinsConfig.LoadConfig()
	if err != nil {
		return errors.New("❌ cannot load config file: " + jenkinsConfig.ConfigFullPath)
	}

	j.JenkinsUser = config.JenkinsUser
	j.Server = config.Server
	j.Token = config.Token
	j.Context = context.Background()

	j.Instance = gojenkins.CreateJenkins(
		nil,
		j.Server,
		j.JenkinsUser,
		j.Token)

	return nil
}

// ServerInfo will show information regarding the server
//
// Args:
//
func (j *Jenkins) ServerInfo() error {
	err := serverReachable(j.Server)
	if err != nil {
		return errors.New("❌ jenkins server unreachable: " + j.Server)
	}

	j.Instance.Info(j.Context)
	fmt.Printf("✅ Connected with: %s\n", j.JenkinsUser)
	fmt.Printf("✅ Server: %s\n", j.Server)
	fmt.Printf("✅ Version: %s\n", j.Instance.Version)

	return nil
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
	return nil

}
