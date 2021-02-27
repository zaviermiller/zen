package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/zaviermiller/zen/internal/vm"
)

func main() {
	cmd := exec.Command(fmt.Sprintf("make release v='%s'", vm.VERSION.String()))
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout

	cmd.Run()
}
