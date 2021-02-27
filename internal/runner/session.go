package runner

import (
	"fmt"
	"time"

	"github.com/zaviermiller/zen/internal/diff"
	d "github.com/zaviermiller/zen/internal/display"
	"github.com/zaviermiller/zen/internal/vm"
)

// implemenation of proc listener for zen session
type ZenListenerImpl struct {
}

func (z ZenListenerImpl) OnInput(zProcess ZenProcess, zSession *ZenSession, i int) {
	if zProcess.Type == TEST {
		d.PrintLoader(i+1, len(zSession.Inputs), "Correct lab terminated. Executing test lab...")
	}
}

func (z ZenListenerImpl) OnComplete(zProcess ZenProcess, zSession ZenSession) {
	zSession.InputNotifier.unregisterAll()
}

// ZenSession is a struct to manage to manage the state of
// ... well, a basic zen session.
type ZenSession struct {
	Inputs         []string
	InputNotifier  InputNotifier
	ProcListener   ZenProcessListener
	CorrectProcess *ZenProcess
	TestProcess    *ZenProcess
	Diff           diff.ZenDiff
}

// NewSession creates and returns a new ZenSession given two ZenProcesses
func NewSession(correct *ZenProcess, test *ZenProcess) (ZenSession, error) {
	notifier := InputNotifier{}
	session := ZenSession{Inputs: []string{}, CorrectProcess: correct, TestProcess: test, InputNotifier: notifier, ProcListener: ZenListenerImpl{}}

	return session, nil
}

// Diff calculates and returns whichever diff is instantiated here.
func (z ZenSession) GetDiffs() (diff.ZenDiff, float64) {
	// line length logic if not the same length
	zDiff := diff.NewDiff(diff.MYERS, z.CorrectProcess.CmdName, z.TestProcess.CmdName)

	score := zDiff.Calculate(z.CorrectProcess.GetResponses(), z.TestProcess.GetResponses())

	return zDiff, score

}

// Run executes a ZenSession, going through both binaries and finding all diffs.
func (z *ZenSession) Run() error {

	d.PrintLogo(vm.VERSION.String())
	d.ZenLogLn("Executing correct lab file (you may need to enter input)...\n")

	err := z.CorrectProcess.Run(z, true)
	if err != nil {
		return err
	}

	d.PrintLoader(0, len(z.Inputs), "Correct lab terminated. Executing test lab...")

	t1 := time.Now()

	err = z.TestProcess.Run(z, false)
	if err != nil {
		return err
	}
	z.TestProcess.Output = z.TestProcess.Output

	fmt.Println("")

	// begin showing comp loader
	// d.PrintLoader(0, max(len(z.CorrectProcess.Output), len(z.TestProcess.Output)), "Finding diffs...")
	d.ZenLog("Finding diffs...")

	// diff vars
	// go through diffs
	// for correctIndex < len(z.CorrectProcess.Output) {
	// 	if testIndex >= len(z.TestProcess.Output) {
	// 		break
	// 	}
	// testOutput := z.CorrectProcess.Output[testIndex]
	// output := z.TestProcess.Output[correctIndex]

	// if output.Content != testOutput.Content {
	// 	// check two lines
	// 	for i := 1; i < 3; i++ {
	// 		if testIndex+i < len(z.TestProcess.Output) {
	// 			if output.Content == z.TestProcess.Output[testIndex+i].Content {
	// 				diff.BasicDiff{Type: "LINE", CorrectDiff: "", TestDiff: testOutput.Content}.Print()
	// 				diffCount++
	// 				testIndex += i
	// 				testOutput = z.TestProcess.Output[testIndex]
	// 				break
	// 			}
	// 		}

	// 		if correctIndex+i < len(z.CorrectProcess.Output) {
	// 			if z.CorrectProcess.Output[i+correctIndex].Content == z.TestProcess.Output[testIndex].Content {
	// 				diff.BasicDiff{Type: "LINE", CorrectDiff: output.Content, TestDiff: ""}.Print()
	// 				diffCount++
	// 				correctIndex += i
	// 				output = z.CorrectProcess.Output[correctIndex]
	// 				break
	// 			}
	// 		}

	// 	}
	// }

	var score float64
	z.Diff, score = z.GetDiffs()

	fmt.Println(d.Bright + "Done ✔" + d.Normal)

	if z.Diff != nil {
		z.Diff.Print()
	}

	// testIndex++
	// correctIndex++
	// time.Sleep(10 * time.Millisecond)
	// d.PrintLoader(max(correctIndex, testIndex), max(len(z.CorrectProcess.Output), len(z.TestProcess.Output)), "Finding diffs...")
	// }

	// get elapsed time
	t2 := time.Now()
	elapsed := t2.Sub(t1)

	d.ZenLogLn(fmt.Sprintf("Finished comparing files in %s", elapsed))

	if score >= 1 {
		fmt.Println(d.Bright + d.Green + "\nNO DIFFERENCES DETECTED! 100% GOOD JOB!\n" + d.Normal)
		return nil
	}

	// show score
	d.ZenLogLn("Score: " + d.ColorScore((score * 100)) + "\n")
	return nil

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
