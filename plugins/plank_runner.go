package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/zaviermiller/zen/pkg"
	"github.com/zaviermiller/zen/pkg/plugins"
)

func init() {
	plugins.ZenPluginRegistry.Runner = &PlankRunner{}
}

type PlankRunner struct {
	BinExecutor pkg.Executor
	GSExecutor  plugins.DefaultExecutor

	GSDir    string
	NumTests int
}

func (r PlankRunner) Args() []pkg.RunnerArg {
	return []pkg.RunnerArg{
		{
			DisplayName: "test binary",
			ArgRequired: true,
		},
		{
			DisplayName: "gradescript dir",
			ArgRequired: true,
		},
	}
}

func (r *PlankRunner) Init(args []string, executor pkg.Executor) error {

	if len(args) < 2 {
		return errors.New("Must have 2 args")
	}

	r.BinExecutor = executor.Clone()
	r.BinExecutor.SetPath(args[0])

	r.GSDir = args[1]
	r.GSExecutor = plugins.DefaultExecutor{Path: r.GSDir + "/gradescript"}

	files, _ := ioutil.ReadDir(r.GSDir + "/Gradescript-Examples")

	r.NumTests = len(files)

	fmt.Println("found ", r.NumTests, " files")

	return nil
}

func (r *PlankRunner) Run(resultsChan chan pkg.Results, rtArgv []string) {
	cwd, _ := os.Getwd()

	os.Mkdir(cwd+"/plank-tmp", 0777)
	defer os.RemoveAll(cwd + "/plank-tmp")

	os.Chdir(cwd + "/plank-tmp")

	r.GSExecutor.Execute(os.Stdin, rtArgv)
}
