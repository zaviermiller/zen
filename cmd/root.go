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
	"os"

	"github.com/zaviermiller/zen/internal/runtime"
	"github.com/zaviermiller/zen/pkg"
	"github.com/zaviermiller/zen/pkg/plugins"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/zaviermiller/zen/internal/config"
)

var cfg config.ZenConfig

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "zen [flags] ...executables",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ArbitraryArgs,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		err := cfg.AppCfg.Runner.Init(args, cfg.AppCfg.Executor)
		cobra.CheckErr(err)

		fmt.Println("Current runner: ", cfg.AppCfg.RunnerID)

		fmt.Printf("%+v", cfg.AppCfg.Runner)

		rt := runtime.Runtime{}

		rt.AddCommand("run", runtime.Command{
			Help: "runs the configured test runner",
			Run: func() {

				resultChan := make(chan pkg.Results)
				defer close(resultChan)

				go cfg.AppCfg.Runner.Run(resultChan, nil)

				result := <-resultChan

				fmt.Println(result)
			},
		})

		for runArgs := rt.ShowInput(); runArgs != nil; runArgs = rt.ShowInput() {
			rt.RunCommand(runArgs)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.Flags().StringP("runner", "r", "", "Override the configured runner.")
	// rootCmd.Flags().StringP("executor", "x", "", "Override the configured executor.")

	if argRunner, ok := cfg.AppCfg.Runner.(pkg.ArgsRunner); ok {
		argArr := argRunner.Args()
		// rootCmd.Args = cobra.ExactArgs(len(argArr))
		rootCmd.Use = "zen [flags] " + RunnerArgsUsage(argArr)

	}

	if flagRunner, ok := cfg.AppCfg.Runner.(pkg.FlagsRunner); ok {
		flags := flagRunner.Flags()
		flagset := pflag.NewFlagSet("runtime flags", pflag.ExitOnError)

		fmt.Println(flags, flagset)
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	// Find home directory.
	home, err := os.UserHomeDir()
	zenDir := home + "/.zen"

	cobra.CheckErr(err)

	// load the config data
	cfg = config.InitConfig(zenDir)

	// set default exec and runner
	cfg.AppCfg.Executor = &plugins.DefaultExecutor{}
	cfg.AppCfg.Runner = &plugins.DefaultRunner{}

	// load plugin from path based on config default
	for _, plugin := range cfg.Plugins.Runners {
		if plugin.Name == cfg.AppCfg.RunnerID {
			plugins.LoadPlugin(plugin.Path)
			cfg.AppCfg.Runner = plugins.ZenPluginRegistry.Runner
		}
	}

	for _, plugin := range cfg.Plugins.Executors {
		if plugin.Name == cfg.AppCfg.ExecutorID {
			plugins.LoadPlugin(plugin.Path)
			cfg.AppCfg.Executor = plugins.ZenPluginRegistry.Executor
		}
	}

}

func RunnerArgsUsage(argDefs []pkg.RunnerArg) string {
	out := ""
	for _, arg := range argDefs {
		if arg.ArgRequired {
			out += arg.DisplayName + " "
		} else {
			out += "[" + arg.DisplayName + "] "
		}
	}

	return out
}
