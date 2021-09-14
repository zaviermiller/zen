package plugins

import (
	"os"

	"github.com/zaviermiller/zen/pkg"
)

type DefaultRunner struct {
	TestExec    pkg.Executor
	CorrectExec pkg.Executor
}

func (r *DefaultRunner) Init(args []string, exec pkg.Executor) error {
	r.TestExec = exec.Clone()
	r.CorrectExec = exec.Clone()

	r.TestExec.SetPath(args[0])
	r.CorrectExec.SetPath(args[1])

	return nil
}

func (r *DefaultRunner) Args() []pkg.RunnerArg {
	return []pkg.RunnerArg{
		{
			DisplayName: "test binary",
			ArgRequired: true,
		},
		{
			DisplayName: "correct binary",
			ArgRequired: false,
		},
	}
}

func (r *DefaultRunner) Run(emitResults chan pkg.Results, rtArgv []string) {
	testOutput, _ := r.TestExec.Execute(os.Stdin, rtArgv)
	correctOutput, _ := r.CorrectExec.Execute(os.Stdin, rtArgv)

	result := pkg.Results{Correct: testOutput.Equals(correctOutput), Expected: testOutput.Stdout, Actual: correctOutput.Stdout}

	// fmt.Println(result)

	emitResults <- result

}

// var DefaultZenRunner DefaultRunner = DefaultRunner{}
