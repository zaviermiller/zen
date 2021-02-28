package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"os/exec"

	d "github.com/zaviermiller/zen/internal/display"
	"github.com/zaviermiller/zen/internal/runner"
	"github.com/zaviermiller/zen/internal/vm"
)

func main() {

	// flags
	versionFlag := flag.Bool("v", false, "show the version number")

	flag.Parse()

	if *versionFlag {
		ShowVersion()
		return
	}

	if vm.CheckUpdate() {
		return
	}

	// diff := diff.NewDiff(diff.MYERS, 0, 0, "correct", "test")
	// diff.Calculate("zavier", "majed")

	// diff.Print()

	// ensure correct usage
	if len(os.Args) < 3 {
		check(errors.New("Usage: " + d.Purple + "zen [example binary] [your binary] [flags]" + d.Normal))
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
		fmt.Println(d.Bright + d.Red + "ZEN ERROR: " + d.Normal + err.Error())
		os.Exit(1)
	}
}

// ShowVersion shows the Zen ascii logo and installed version
func ShowVersion() {
	fmt.Println(d.Bright + d.Purple + d.Zen(vm.VERSION.String()) + "\n")
}
