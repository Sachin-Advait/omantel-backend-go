package config

import (
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gitlab.com/rendezvousPoint/incentive_calculator/commons"
)

func NewViperConfig() *Configuration {
	viper.SetConfigName(commons.CONFIG_PATH)
	err := viper.BindEnv("gopath")
	if err != nil {
		logrus.WithError(err).Fatal("gopath load failed")
	}
	viper.AddConfigPath(viper.GetString("gopath") + "/src/gitlab.com/rendezvousPoint/incentive_calculator/config")
	//viper.AddConfigPath(commons.CONFIG_DIR)

	if err = viper.ReadInConfig(); err != nil {
		logrus.Fatalf("Error reading config file, %s", err)
	}

	var configuration Configuration
	err = viper.Unmarshal(&configuration)
	if err != nil {
		logrus.Fatalf("unable to decode into struct, %v", err)
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		err := viper.Unmarshal(&configuration)
		if err != nil {
			logrus.Fatalf("unable to decode into struct, %v", err)
		}
		logrus.Info("Config file changed:", e.Name)
	})

	return &configuration
}

var once sync.Once
var conf *Configuration

func GetInstance() *Configuration {
	once.Do(func() {
		conf = NewViperConfig()
		defer logrus.Info("config setup completed")
	})

	return conf
}
