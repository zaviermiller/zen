package plugins

import (
	"plugin"

	"github.com/zaviermiller/zen/pkg"
)

type PluginType int

const (
	RunnerType   PluginType = 1
	ExecutorType PluginType = 2
	LastLoaded   PluginType = 128
)

type PluginRegistry struct {
	Executor pkg.Executor
	Runner   pkg.Runner
}

var ZenPluginRegistry = PluginRegistry{
	Executor: nil,
	Runner:   nil,
}

func LoadPlugin(path string) error {
	if _, err := plugin.Open(path); err != nil {
		return err
	}

	return nil
}

func Load(paths ...string) error {
	for _, path := range paths {
		if _, err := plugin.Open(path); err != nil {
			return err
		}
	}
	return nil
}
