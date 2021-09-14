package plugins

import (
	. "github.com/zaviermiller/zen/pkg"
)

type ExecutorPlugin interface {
	BuildExecutor(options map[string]interface{}) Executor
}

type flagged struct {
	Runner
	flags map[string]RunnerFlag
}

func (f *flagged) Flags() map[string]RunnerFlag {
	return f.flags
}

func RegisterRunnerFlags(runner Runner, flags map[string]RunnerFlag) Runner {
	return &flagged{
		runner,
		flags,
	}
}

type argged struct {
	Runner
	args []RunnerArg
}

func (a *argged) Args() []RunnerArg {
	return a.args
}

func RegisterRunnerArgs(runner Runner, args []RunnerArg) Runner {
	return &argged{
		runner,
		args,
	}
}

/*
Plugins:
	The point of plugins is to have runtime configurable behaviour that can
	be installed by the end user without having to edit any code.

	Currently, I can see three distinct plugin types: RunnerPlugin, ExecutorPlugin, and RuntimePlugin.
	The details for each are as follows:

		RunnerPlugin:
			Creates a runner that can be used with a configured executor (or depend on a manual one?)
			The runner can define flags and the positional arguments that are required to run the tests.
			For example:
									`zen -r RunnerPluginImpl test-to-exec correct-to-exec` <- this is interpreted correctly with the runner, and will show
																																								in the help message if set in config. Each arg can also have a defined executor
									`> run [flag options]` <- flags should probably be passed into the run function as arguments, so well need a way to define flags
									`> customcmd [flag options] [other args]` <- this is more a runtime plugin feature, but custom cmd support as well (maybe runner is also runtime plugin)

		ExecutorPlugin:
			This simply defines a configured executor that will be tried with the runner. Executors can be configured by runners manually,
			But when not manually overridden, it can be set indendently.

		RuntimePlugin:
			Defines new commands that can be run, going to try and figure this out later.

*/
