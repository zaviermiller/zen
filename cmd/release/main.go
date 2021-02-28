package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zaviermiller/zen/internal/vm"
)

func main() {
	cmd := exec.Command("make", "compile", fmt.Sprintf("v='%s'", vm.VERSION.FileString()))
	cmd.Stderr = os.Stderr
	// cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Run()
}
