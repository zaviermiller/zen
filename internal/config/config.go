package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"

	. "github.com/zaviermiller/zen/pkg"
)

type ZenConfig struct {
	AppCfg  AppConfig `mapstructure:"app"`
	Plugins PluginsConfig
}

func (z ZenConfig) Write() error {
	if err := viper.WriteConfig(); err != nil {
		return err
	}
	return nil
}

func InitConfig(cfgDir string) ZenConfig {
	var cfg ZenConfig

	if _, err := os.Stat(cfgDir); os.IsNotExist(err) {
		// First load, create config and data files
		if err := os.Mkdir(cfgDir, 0777); err != nil {
			panic("couldnt make directory at " + cfgDir)
		}

	} else {
		if _, err := os.OpenFile(cfgDir+".plugins.toml", os.O_CREATE, 0644); err != nil {
			panic("couldnt open plugins file")
		}
	}

	viper.AddConfigPath(cfgDir)
	viper.SetConfigType("toml")
	viper.SetConfigName("config")

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())

		if err = viper.Unmarshal(&cfg); err != nil {
			fmt.Fprintln(os.Stderr, "Couldn't unmarshal config into struct.")
		}
	} else {
		// create file if not found
		err = viper.SafeWriteConfig()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Could not write config file to ", viper.ConfigFileUsed())
		}
	}

	// find the plugin?

	return cfg
}

type AppConfig struct {
	ExecutorID string `mapstructure:"executor"`
	RunnerID   string `mapstructure:"runner"`

	Executor Executor `mapstructure:"-"`
	Runner   Runner   `mapstructure:"-"`
}

type PluginsConfig struct {
	Executors []PluginConfig
	Runners   []PluginConfig
}

// todo: import from plugin
type PluginConfig struct {
	Name string
	Path string
}
