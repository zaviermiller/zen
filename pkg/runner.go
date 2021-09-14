package pkg

type Results struct {
	Correct  bool
	Expected string
	Actual   string
	// other stuff tbd
}

/*
	Runner should have runtime plugin?
	Ok so need a way to define flags, bin args and whether they have defined executor
*/

type Runner interface {
	// Usage returns the message to display when the
	// Runner is default and the zen cmd is run incorrectly
	// Init(args []string) error

	Init(args []string, executor Executor) error

	// Help returns the available options for the Run functions,
	// for example, flags
	// Help() string

	// Run runs the tests
	Run(emitTestResults chan Results, rtArgv []string)

	// DisplayName() string
}

type RunnerArg struct {
	DisplayName      string
	RequiredExecutor *Executor
	ArgRequired      bool
}

type RunnerFlag struct {
	Name         string
	Shorthand    [3]byte
	Help         string
	DefaultValue interface{}
}

type ArgsRunner interface {
	Args() []*RunnerArg
}

type FlagsRunner interface {
	Flags() []*RunnerFlag
}
