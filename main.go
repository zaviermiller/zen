package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"

	"github.com/zaviermiller/zen/runner"
	u "github.com/zaviermiller/zen/util"
	"github.com/zaviermiller/zen/vm"
)

func main() {

	if vm.CheckUpdate() {
		return
	}

	// ensure correct usage
	if len(os.Args) < 3 {
		check(errors.New("Usage: " + u.Purple + "zen [example binary] [your binary]" + u.Normal))
	}

	// create variables from args
	correctBinPath := os.Args[1]
	testBinPath := os.Args[2]
	binOpts := os.Args[3:]

	// build commands for processes
	correctCmd := exec.Command(correctBinPath, binOpts...)
	testCmd := exec.Command(testBinPath, binOpts...)

	// create processes
	correctProc, err := runner.NewProcess(correctCmd, runner.CORRECT)
	check(err)
	testProc, err := runner.NewProcess(testCmd, runner.TEST)
	check(err)

	// create the zen runner
	session, err := runner.NewSession(&correctProc, &testProc)
	check(err)

	// run the zen session
	session.Run()

}

// basic error checking function
func check(err error) {
	if err != nil {
		fmt.Println(u.Bright + u.Red + "ZEN ERROR: " + u.Normal + err.Error())
		os.Exit(1)
	}
}
