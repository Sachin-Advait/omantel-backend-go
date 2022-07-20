package configs

import (
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/spf13/viper"
)

type Clients struct {
	Client
	env map[string]string
}

// Client is the instance of the config client to be used by the application
var client *Clients

type providers struct {
	providers map[string]*viper.Viper
	mu        sync.Mutex
}

var baseConfigPath string
var p *providers

// Init is used to initialize the configurations
func Init(path string) {
	baseConfigPath = path
	p = &providers{
		providers: make(map[string]*viper.Viper),
	}
}

// Get is used to get the instance to the config provider for the configuration name
func Get(name string) (*viper.Viper, error) {
	p.mu.Lock()
	defer p.mu.Unlock()
	// see for an existing provider
	if provider, ok := p.providers[name]; ok {
		// provider already exists
		// use it
		return provider, nil
	}

	// try to get the provider
	provider := viper.New()
	provider.SetConfigName(name)
	provider.AddConfigPath(baseConfigPath)
	err := provider.ReadInConfig()
	if err != nil {
		// config not found
		return nil, fmt.Errorf("config %s not found", name)
	}

	// add a watcher for this provider
	provider.WatchConfig()

	// successfully found config, store it for future use
	p.providers[name] = provider

	return provider, nil
}

// Get is used to get the instance of the client
func GetConfig() *Clients {
	return client
}

// InitTestModeConfigs is used to initialize the configs
func InitTestModeConfigs(directory string, configNames ...string) error {
	c, err := New(Options{
		Provider: FileBased,
		Params: map[string]interface{}{
			"configsDirectory": directory,
			"configNames":      configNames,
			"configType":       "json",
		},
	})
	if err != nil {
		return err
	}
	client = getClient(c)
	return nil
}

// GetStringWithEnv is used to get the config by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *Clients) GetStringWithEnv(config, key string) (string, error) {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return s, err
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s, nil
}

// GetStringWithEnvD is used to get the config with default value by filling the variables from environment variables
// for example, say a config value is ${XYZ}/abc, and the value of environment variable XYZ is ABC,
// then this function will return XYZ/abc.
func (c *Clients) GetStringWithEnvD(config, key, defaultValue string) string {
	// first fetch the config value
	s, err := c.GetString(config, key)
	// if error no pointing moving ahead
	if err != nil {
		return defaultValue
	}
	// now time to look for and replace with all the environment variables
	for k, v := range c.env {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%s}", k), v)
	}
	return s
}
func getEnvironment() map[string]string {
	env := os.Environ()
	result := make(map[string]string)
	for _, e := range env {
		s := strings.Split(e, "=")
		if len(s) == 2 {
			result[s[0]] = s[1]
		}
	}
	return result
}

func getClient(c Client) *Clients {
	return &Clients{Client: c, env: getEnvironment()}
}
