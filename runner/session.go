package runner

import (
	"fmt"
	"time"

	"github.com/zaviermiller/zen/diff"
	u "github.com/zaviermiller/zen/util"
	"github.com/zaviermiller/zen/vm"
)

// implemenation of proc listener for zen session
type ZenListenerImpl struct {
}

func (z ZenListenerImpl) OnInput(zProcess ZenProcess, zSession *ZenSession, i int) {
	if zProcess.Type == TEST {
		u.PrintLoader(i+1, len(zSession.Inputs), "Correct lab terminated. Executing test lab...")
	}
}

func (z ZenListenerImpl) OnComplete(zProcess ZenProcess, zSession ZenSession) {
	zSession.InputNotifier.unregisterAll()
}

type ZenSession struct {
	Inputs         []string
	InputNotifier  InputNotifier
	ProcListener   ZenProcessListenerIface
	CorrectOutputs []ZenOutput
	TestOutputs    []ZenOutput
	CorrectProcess *ZenProcess
	TestProcess    *ZenProcess
	Diffs          []*diff.ZenDiff
}

func NewSession(correct *ZenProcess, test *ZenProcess) (ZenSession, error) {
	notifier := InputNotifier{}
	session := ZenSession{Inputs: []string{}, CorrectOutputs: []ZenOutput{}, TestOutputs: []ZenOutput{}, CorrectProcess: correct, TestProcess: test, InputNotifier: notifier, ProcListener: ZenListenerImpl{}}

	return session, nil
}

func (z *ZenSession) Run() error {

	u.PrintLogo(vm.VERSION.String())
	u.ZenLogLn("Executing correct lab file (you may need to enter input)...\n")

	err := z.CorrectProcess.Run(z, true)
	if err != nil {
		return err
	}
	z.CorrectOutputs = z.CorrectProcess.Output

	fmt.Println("")
	u.PrintLoader(0, len(z.Inputs), "Execution finished! Executing test lab...")

	t1 := time.Now()

	err = z.TestProcess.Run(z, false)
	if err != nil {
		return err
	}
	z.TestOutputs = z.TestProcess.Output

	fmt.Println("")

	// begin showing comp loader
	u.PrintLoader(0, max(len(z.CorrectOutputs), len(z.TestOutputs)), "Finding diffs...")

	// diff vars
	diffCount := 0
	testIndex := 0
	correctIndex := 0

	// go through diffs
	for correctIndex < len(z.CorrectOutputs) {
		if testIndex >= len(z.TestOutputs) {
			break
		}
		testOutput := z.CorrectOutputs[testIndex]
		output := z.TestOutputs[correctIndex]

		if output.Content != testOutput.Content {
			// check two lines
			for i := 1; i < 3; i++ {
				if testIndex+i < len(z.TestOutputs) {
					if output.Content == z.TestOutputs[testIndex+i].Content {
						diff.ZenDiff{Type: "LINE", CorrectDiff: "", TestDiff: testOutput.Content}.Print(correctIndex, testIndex, z.CorrectProcess.CmdName, z.TestProcess.CmdName)
						diffCount++
						testIndex += i
						testOutput = z.TestOutputs[testIndex]
						break
					}
				}

				if correctIndex+i < len(z.CorrectOutputs) {
					if z.CorrectOutputs[i+correctIndex].Content == z.TestOutputs[testIndex].Content {
						diff.ZenDiff{Type: "LINE", CorrectDiff: output.Content, TestDiff: ""}.Print(correctIndex, testIndex, z.CorrectProcess.CmdName, z.TestProcess.CmdName)
						diffCount++
						correctIndex += i
						output = z.CorrectOutputs[correctIndex]
						break
					}
				}

			}
		}

		diff := output.CompareDiff(testOutput)

		if diff != nil {
			diffCount++
			diff.Print(correctIndex, testIndex, z.CorrectProcess.CmdName, z.CorrectProcess.CmdName)
		}

		testIndex++
		correctIndex++
		time.Sleep(10 * time.Millisecond)
		u.PrintLoader(max(correctIndex, testIndex), max(len(z.CorrectOutputs), len(z.TestOutputs)), "Finding diffs...")
	}

	// get elapsed time
	t2 := time.Now()
	elapsed := t2.Sub(t1)

	u.ZenLogLn(fmt.Sprintf("Finished comparing files in %s", elapsed))

	if diffCount == 0 {
		fmt.Println(u.Bright + u.Green + "\nNO DIFFERENCES DETECTED! 100% GOOD JOB!\n" + u.Normal)
		return nil
	}

	// show score
	u.ZenLogLn("Score: " + u.ColorScore(float64(len(z.CorrectOutputs)-diffCount)/float64(len(z.CorrectOutputs))) + "\n")
	return nil

}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
