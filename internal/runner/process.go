package runner

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

type ProcessType int

const (
	CORRECT ProcessType = iota
	TEST
)

type ZenProcess struct {
	CmdName       string
	Type          ProcessType
	Args          []string
	Process       *os.Process
	Cmd           *exec.Cmd
	Output        []ZenOutput
	Err           error
	InputDetector *InputDetector
	InPipe        io.WriteCloser
	OutPipe       io.ReadCloser
}

type ZenProcessListener interface {
	OnComplete(zProcess ZenProcess, zSession ZenSession)
	// OnError(zProcess ZenProcess, err error)
	OnInput(zProcess ZenProcess, zSession *ZenSession, i int)
}

func NewProcess(cmd *exec.Cmd, t ProcessType) (ZenProcess, error) {
	execIn, err := cmd.StdinPipe()
	if err != nil {
		return ZenProcess{}, err
	}
	execOut, err := cmd.StdoutPipe()
	if err != nil {
		return ZenProcess{}, err
	}

	detector := InputDetector{}

	proc := ZenProcess{CmdName: os.Args[0], Type: t, Args: os.Args[1:], Cmd: cmd, InPipe: execIn, OutPipe: execOut, InputDetector: &detector}

	return proc, nil
}

func (z *ZenProcess) Run(zs *ZenSession, interactable bool) error {

	// create the output scanner and scan every byte (char)
	scanner := bufio.NewScanner(z.OutPipe)
	scanner.Split(bufio.ScanBytes)

	// async function to read the output
	go func() {

		// keep track of output
		tmpString := ""
		if interactable {
			fmt.Print("   ")
		}

		// create and register input detector to break off output of type PROMPT on input
		z.InputDetector.update = func() {
			z.Output = append(z.Output, ZenOutput{Type: PROMPT, Content: tmpString})

			// reset
			tmpString = ""
			if interactable {
				fmt.Print("   ")
			}
		}
		zs.InputNotifier.register(*z.InputDetector)

		// begin scanning output
		for scanner.Scan() {
			text := scanner.Text()
			tmpString += text
			if interactable {
				fmt.Print(text)
			}

			// break and reset output of type RESPONSE when new line
			if text == "\n" {
				z.Output = append(z.Output, ZenOutput{Type: RESPONSE, Content: tmpString})
				tmpString = ""
				if interactable {
					fmt.Print("   ")
				}
			}
		}
	}()

	err := z.Cmd.Start()
	if err != nil {
		return err
	}

	if interactable {
		go func() {

			// stdin reader
			reader := bufio.NewReader(os.Stdin)

			// read on 'enter'
			input, err := reader.ReadString('\n')
			for err == nil {
				// notify input
				zs.InputNotifier.notifyInput()

				// add input to input list
				zs.Inputs = append(zs.Inputs, input)
				n, err := io.WriteString(z.InPipe, input)
				if err != nil || n == 0 {
					// err means no more input needed
					break
				}
				input, err = reader.ReadString('\n')
			}

			// close input pipe
			z.InPipe.Close()
		}()
	} else {
		if len(zs.Inputs) < 1 {
			return errors.New("No inputs for automated process")
		}

		// do each input with time between
		for i, input := range zs.Inputs {
			time.Sleep(300 * time.Millisecond)
			zs.InputNotifier.notifyInput()
			io.WriteString(z.InPipe, input)
			zs.ProcListener.OnInput(*z, zs, i)
		}
	}

	err = z.Cmd.Wait()
	if err != nil {
		return err
	}

	zs.ProcListener.OnComplete(*z, *zs)
	return nil

}

// GetResponses returns all the detected response outputs
func (z ZenProcess) GetResponses() []string {
	tmp := []string{}
	for _, out := range z.Output {
		if out.Type == RESPONSE {
			tmp = append(tmp, out.Text())
		}
	}

	return tmp
}
