package runtime

import (
	"fmt"
	"strings"
)

type Runtime struct {
	Commands map[string]Command
}

type Command struct {
	Help      string
	Run       func()
	Shorthand string
}

func (r *Runtime) AddCommand(name string, cmd Command) {
	if r.Commands == nil {
		r.Commands = map[string]Command{}
	}
	r.Commands[name] = cmd
}

func (r *Runtime) ShowInput() []string {
	var inp string

	fmt.Print("> ")
	fmt.Scanln(&inp)

	args := strings.Split(inp, " ")

	if args[0] == "quit" {
		return nil
	} else if args[0] == "help" {
		return args
	} else if _, ok := r.Commands[args[0]]; !ok {
		return nil
	}

	return args

}

func (r *Runtime) RunCommand(args []string) {
	if args[0] == "help" {
		for name, cmd := range r.Commands {
			fmt.Printf("%10s -- %s\n", name, cmd.Help)
		}
	} else {
		cmd := r.Commands[args[0]]

		cmd.Run()
	}

}
