package plugins

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"

	"github.com/zaviermiller/zen/pkg"
	zen "github.com/zaviermiller/zen/pkg"
)

type DefaultExecutor struct {
	Path string
}

func (e *DefaultExecutor) SetPath(path string) {
	e.Path = path
}

func (e *DefaultExecutor) Clone() pkg.Executor {
	newExec := *e

	return &newExec
}

func (e *DefaultExecutor) Execute(stdin io.Reader, argv []string) (zen.Output, error) {
	var err, out bytes.Buffer

	fmt.Println("Beginning exec")
	cmd := exec.Command(e.Path, argv...)

	// set the commands stdin to the passed reader
	cmd.Stdin = stdin

	cmd.Stdout = &out
	cmd.Stderr = &err

	cmd.Run()

	output := zen.Output{Stdout: out.String(), Stderr: err.String()}

	return output, nil
}
