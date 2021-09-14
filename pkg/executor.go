package pkg

import "io"

// Executors are responsible for correctly running a program
// and returning its stdout and stderr in an Output object.

type Executor interface {
	Execute(path string, stdin io.Reader, argv []string) (Output, error)
}

// type InputRunner interface {
// 	Exec(path string, argv []string, stdin []string) Output
// }
