package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	u "github.com/zaviermiller/zen/util"
)

// quick and dirty observer implementation
// to know when an input causes newline

type inputDetector struct {
	update func()
}

type inputSubject struct {
	observers []inputDetector
}

// register a new detector to listen to the subject
func (s *inputSubject) register(i inputDetector) {
	s.observers = append(s.observers, i)
}

// remove all listeners (easy to implement and all i needed really)
func (s *inputSubject) unregisterAll() {
	s.observers = []inputDetector{}
}

// notify all detectors of a change
func (s inputSubject) notifyInput() {
	for _, detector := range s.observers {
		detector.detected()
	}
}

// the detected func runs the passed update func from initialization,
// in order to have access to necessary data
func (d *inputDetector) detected() {
	d.update()
}

// enum type for different outputs
type outputType int

const (
	PROMPT outputType = iota
	RESPONSE
)

// output type that can differentiate between prompt and response outputs
type output struct {
	outputType outputType
	content    string
}

// comparison of two outputs returning a diff object
func (o output) compareDiff(o1 output) *zenDiff {
	// display constant
	om := []string{"PROMPT", "RESPONSE"}

	// check to make sure the output types are the same (worst case if outputTypes dont match)
	if o.outputType != o1.outputType {
		return &zenDiff{errorType: "OUTPUT TYPE", testDiff: om[o1.outputType], correctDiff: om[o.outputType]}
	}

	// basic forward checking of output to ensure match
	// probably need to improve
	for i, char := range o.content {
		if char != rune(o1.content[i]) {
			return &zenDiff{errorType: om[o1.outputType] + " CONTENT", testDiff: o1.content[i:], correctDiff: o.content[i:]}
		}
	}

	// lines matched
	return nil
}

// difference object that contains data for printing
// the diffs to the screen
type zenDiff struct {
	testDiff    string
	correctDiff string
	errorType   string
}

// prints out the diff
func (z zenDiff) print(correctId, testId int, correctBinPath, testBinPath string) {
	fmt.Println(u.Red + u.Bright + u.Clear + "!!!" + z.errorType + " DIFF!!! " + u.Normal + "\n")
	fmt.Println("   " + correctBinPath + " [" + strconv.Itoa(correctId) + "]: " + u.Green + z.correctDiff + u.Normal)
	fmt.Println("   " + testBinPath + " [" + strconv.Itoa(testId) + "]: " + u.Red + z.testDiff + u.Normal + "\n")
}

func main() {

	// ensure correct usage
	if len(os.Args) < 3 {
		check(errors.New("Usage: " + u.Purple + "zen [example binary] [your binary]" + u.Normal))
	}

	// create variables from args
	correctBinPath := os.Args[1]
	testBinPath := os.Args[2]
	binOpts := os.Args[3:]

	// inputs and correct outputs store
	inputs := []string{}
	outputs := []output{}

	// input notifier
	inpSubject := inputSubject{}

	// pretty stuff
	fmt.Println(u.Green + u.Bright + u.Zen)
	fmt.Println(u.Purple + "\n   \"Bring yourself peace\" -Z" + u.Normal + "\n")
	zenLogLn("Executing correct lab file (you may need to enter input)...\n")

	// build correct lab cmd and io pipes
	cmd := exec.Command(correctBinPath, binOpts...)
	cmd.Stderr = os.Stderr

	execIn, err := cmd.StdinPipe()
	check(err)
	execOut, err := cmd.StdoutPipe()
	check(err)

	// create the output scanner and scan every byte (char)
	scanner := bufio.NewScanner(execOut)
	scanner.Split(bufio.ScanBytes)

	// async function to read the output
	go func() {

		// keep track of output
		tmpString := ""
		fmt.Print("   ")

		// create and register input detector to break off output of type PROMPT on input
		inputSignal := inputDetector{update: func() {
			outputs = append(outputs, output{outputType: PROMPT, content: tmpString})

			// reset
			tmpString = ""
			fmt.Print("   ")
		}}
		inpSubject.register(inputSignal)

		// begin scanning output
		for scanner.Scan() {
			text := scanner.Text()
			tmpString += text
			fmt.Print(text)

			// break and reset output of type RESPONSE when new line
			if text == "\n" {
				outputs = append(outputs, output{outputType: RESPONSE, content: tmpString})
				tmpString = ""
				fmt.Print("   ")
			}
		}
	}()

	// start running cmd
	err = cmd.Start()
	check(err)

	// async function to read input
	go func() {

		// stdin reader
		reader := bufio.NewReader(os.Stdin)

		// read on 'enter'
		input, err := reader.ReadString('\n')
		for err == nil {
			// notify input
			inpSubject.notifyInput()

			// add input to input list
			inputs = append(inputs, input)
			n, err := io.WriteString(execIn, input)
			if err != nil || n == 0 {
				// err means no more input needed
				break
			}
			input, err = reader.ReadString('\n')
		}

		// close input pipe
		execIn.Close()
	}()

	// wait for cmd to finish
	cmd.Wait()

	// reset notifier
	inpSubject.unregisterAll()

	// track test exec time
	t1 := time.Now()

	// Start test execution
	zenLog("Execution finished! Executing test lab...")

	// outputs of test execution
	testOutputs := []output{}

	// build test lab cmd and pipes
	cmd = exec.Command(testBinPath, binOpts...)
	testIn, err := cmd.StdinPipe()
	check(err)
	testOut, err := cmd.StdoutPipe()
	check(err)

	// build new scanner for output pipe
	scanner = bufio.NewScanner(testOut)
	scanner.Split(bufio.ScanBytes)

	// async output, same as above
	go func() {
		tmpString := ""
		inputSignal := inputDetector{update: func() {
			testOutputs = append(testOutputs, output{outputType: PROMPT, content: tmpString})
			tmpString = ""
		}}
		inpSubject.register(inputSignal)
		for scanner.Scan() {
			text := scanner.Text()
			tmpString += text

			if text == "\n" {
				testOutputs = append(testOutputs, output{outputType: RESPONSE, content: tmpString})
				tmpString = ""
			}
		}
	}()

	// run test lab
	err = cmd.Start()
	check(err)

	// do each input with time between
	for _, input := range inputs {
		time.Sleep(300 * time.Millisecond)
		inpSubject.notifyInput()
		io.WriteString(testIn, input)
	}

	// wait for lab test to finish
	cmd.Wait()

	// show test exec is done
	fmt.Println(u.Bright + "Done!" + u.Normal + "\n")

	// begin showing comp loader
	printLoader(0, max(len(testOutputs), len(outputs)))

	// diff vars
	diffCount := 0
	testIndex := 0
	correctIndex := 0

	// go through diffs
	for correctIndex < len(outputs) {
		if testIndex >= len(testOutputs) {
			break
		}
		testOutput := testOutputs[testIndex]
		output := outputs[correctIndex]

		if output.content != testOutput.content {
			// check two lines
			for i := 1; i < 3; i++ {
				if testIndex+i < len(testOutputs) {
					if output.content == testOutputs[testIndex+i].content {
						zenDiff{errorType: "LINE", correctDiff: "", testDiff: testOutput.content}.print(correctIndex, testIndex, correctBinPath, testBinPath)
						diffCount++
						testIndex += i
						testOutput = testOutputs[testIndex]
						break
					}
				}

				if correctIndex+i < len(outputs) {
					if outputs[i+correctIndex].content == testOutputs[testIndex].content {
						zenDiff{errorType: "LINE", correctDiff: output.content, testDiff: ""}.print(correctIndex, testIndex, correctBinPath, testBinPath)
						diffCount++
						correctIndex += i
						output = outputs[correctIndex]
						break
					}
				}

			}
		}

		diff := output.compareDiff(testOutput)

		if diff != nil {
			diffCount++
			diff.print(correctIndex, testIndex, correctBinPath, testBinPath)
		}

		testIndex++
		correctIndex++
		time.Sleep(10 * time.Millisecond)
		printLoader(max(correctIndex, testIndex), max(len(testOutputs), len(outputs)))
	}

	// get elapsed time
	t2 := time.Now()
	elapsed := t2.Sub(t1)

	fmt.Println()
	zenLogLn(fmt.Sprintf("Finished comparing files in %s", elapsed))

	if diffCount == 0 {
		fmt.Println(u.Bright + u.Green + "\nNO DIFFERENCES DETECTED! 100% GOOD JOB!\n" + u.Normal)
		return
	}

	// show score
	zenLogLn("Score: " + colorScore(float64(len(outputs)-diffCount)/float64(len(outputs))) + "\n")

}

// basic error checking function
func check(err error) {
	if err != nil {
		fmt.Println(u.Bright + u.Red + "ZEN ERROR: " + u.Normal + err.Error())
		os.Exit(1)
	}
}

// max of two nums
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// prints loader
func printLoader(done, total int) {
	width := 50.0

	progress := float64(done) / float64(total)
	equals := strings.Repeat("=", int(progress*width))
	dashes := strings.Repeat("-", int((1.0-progress)*width))
	fmt.Print(fmt.Sprintf("["+u.Purple+"*"+u.Normal+"] Finding diffs... [ "+equals+dashes+" ] [%d/%d compared] - (%f", int(done), int(total), math.Round(progress*100.0*100.0)/100.0) + "%) \r")
}

// computes color of score
func colorScore(score float64) string {
	textScore := fmt.Sprintf("%.2f", score)
	if score >= .75 {
		return u.Green + textScore + u.Normal
	} else if score >= .5 {
		return u.Yellow + textScore + u.Normal
	} else {
		return u.Red + textScore + u.Normal
	}
}

// log helpers
func zenLog(msg string) {
	fmt.Print("\n\r[" + u.Purple + "*" + u.Normal + "] " + msg)
}

func zenLogLn(msg string) {
	fmt.Println("\n\r[" + u.Purple + "*" + u.Normal + "] " + msg)
}
