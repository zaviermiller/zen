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
	value  bool
	update func()
}

type inputSubject struct {
	observers []inputDetector
}

func (s *inputSubject) register(i inputDetector) {
	s.observers = append(s.observers, i)
}

func (s *inputSubject) unregisterAll() {
	s.observers = []inputDetector{}
}

func (s inputSubject) notifyInput() {
	for _, detector := range s.observers {
		detector.detected()
	}
}

func (d *inputDetector) detected() {
	d.value = true
	d.update()
}

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

func (o output) compareDiff(o1 output) *zenDiff {
	om := []string{"PROMPT", "RESPONSE"}
	if o.outputType != o1.outputType {
		return &zenDiff{errorType: "OUTPUT TYPE ", testDiff: om[o1.outputType], correctDiff: om[o.outputType]}
	}

	// fDiffChan := make(chan int)

	// go func() {
	// 	for i, char := range o.content {
	// 		if char != rune(o1.content[i]) {
	// 			// return &zenDiff{errorType: om[o1.outputType] + " -- CONTENT DIFF", testDiff: o1.content[i:], correctDiff:  }
	// 			fDiffChan <- i
	// 		}
	// 	}
	// }()

	// bDiff := 0

	for i, char := range o.content {
		if char != rune(o1.content[i]) {
			// bDiff := i
			// break
			return &zenDiff{errorType: om[o1.outputType] + " CONTENT ", testDiff: o1.content[i:], correctDiff: o.content[i:]}
		}
	}

	// fDiff := <-fDiffChan

	// if
	return nil
}

type zenDiff struct {
	testDiff    string
	correctDiff string
	errorType   string
}

func (z zenDiff) print(id int, correctBinPath, testBinPath string) {
	fmt.Println("!!!" + u.Red + u.Bright + z.errorType + "DIFF!!! " + u.Normal + "\n")
	fmt.Println("   " + correctBinPath + " [" + strconv.Itoa(id) + "]: " + u.Green + u.Bright + z.correctDiff + u.Normal)
	fmt.Println("   " + testBinPath + " [" + strconv.Itoa(id) + "]: " + u.Red + z.testDiff + u.Normal + "\n")
}

func main() {

	if len(os.Args) < 3 {
		check(errors.New("Usage: " + u.Purple + "zen [example binary] [your binary]" + u.Normal))
	}
	correctBinPath := os.Args[1]
	testBinPath := os.Args[2]
	binOpts := os.Args[3:]

	inputs := []string{}
	outputs := []output{}

	inpSubject := inputSubject{}

	// writer := bufio.NewWriter(os.Stdout)

	fmt.Println(u.Green + u.Bright + u.Zen)
	fmt.Println(u.Purple + "\n   \"Bring yourself peace\" -Z" + u.Normal + "\n")
	fmt.Println("[" + u.Purple + "*" + u.Normal + "] Executing correct lab file (you may need to enter input)...\n")

	cmd := exec.Command(correctBinPath, binOpts...)
	cmd.Stderr = os.Stderr
	// cmd.Stdout = os.Stdout
	// cmd.Stdin = os.Stdin
	execIn, err := cmd.StdinPipe()
	check(err)
	execOut, err := cmd.StdoutPipe()
	check(err)

	scanner := bufio.NewScanner(execOut)
	scanner.Split(bufio.ScanBytes)

	go func() {
		tmpString := ""
		fmt.Print("   ")
		inputSignal := inputDetector{value: false, update: func() {
			outputs = append(outputs, output{outputType: PROMPT, content: tmpString})
			tmpString = ""
			fmt.Print("   ")
		}}
		inpSubject.register(inputSignal)
		for scanner.Scan() {
			text := scanner.Text()
			tmpString += text
			fmt.Print(text)

			if text == "\n" {
				outputs = append(outputs, output{outputType: RESPONSE, content: tmpString})
				tmpString = ""
				fmt.Print("   ")
			}
		}
	}()

	err = cmd.Start()
	check(err)

	go func() {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		for err == nil {
			inpSubject.notifyInput()
			inputs = append(inputs, input)
			n, err := io.WriteString(execIn, input)
			if err != nil || n == 0 {
				break
			}
			input, err = reader.ReadString('\n')
		}
		execIn.Close()
	}()

	cmd.Wait()

	inpSubject.unregisterAll()

	t1 := time.Now()

	// Start test execution
	fmt.Print("\r[" + u.Purple + "*" + u.Normal + "] Execution finished! Executing test lab...")

	testOutputs := []output{}

	cmd = exec.Command(testBinPath, binOpts...)
	testIn, err := cmd.StdinPipe()
	check(err)
	testOut, err := cmd.StdoutPipe()
	check(err)

	scanner = bufio.NewScanner(testOut)
	scanner.Split(bufio.ScanBytes)

	go func() {
		tmpString := ""
		inputSignal := inputDetector{value: false, update: func() {
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

	err = cmd.Start()
	check(err)

	for _, input := range inputs {
		time.Sleep(300 * time.Millisecond)
		inpSubject.notifyInput()
		io.WriteString(testIn, input)
	}

	cmd.Wait()

	fmt.Println(u.Bright + "Done!" + u.Normal + "\n")

	printLoader(0, max(len(testOutputs), len(outputs)))

	diffCount := 0
	testIndex := 0
	correctIndex := 0

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
						testIndex += i
						testOutput = testOutputs[testIndex]
						break
					}
				}

				if correctIndex+i < len(outputs) {
					if outputs[i+correctIndex].content == testOutputs[testIndex].content {
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
			diff.print(correctIndex, correctBinPath, testBinPath)
		}

		testIndex++
		correctIndex++
		// time.Sleep(100 * time.Millisecond)
		printLoader(max(correctIndex, testIndex), max(len(testOutputs), len(outputs)))

		// loaderProgress <- max(testIndex, correctIndex)
	}

	t2 := time.Now()
	elapsed := t2.Sub(t1)

	fmt.Println(fmt.Sprintf("\n\n["+u.Purple+"*"+u.Normal+"] Finished comparing files in %s", elapsed))

	if diffCount == 0 {
		fmt.Println(u.Bright + u.Green + "\nNO MAJOR DIFFERENCES DETECTED! 100% GOOD JOB!\n" + u.Normal)
		return
	}

	fmt.Println("Score: ", float64(diffCount)/float64(len(outputs)))

}

func check(err error) {
	if err != nil {
		fmt.Println(u.Bright + u.Red + "ZEN ERROR: " + u.Normal + err.Error())
		os.Exit(1)
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func printLoader(done, total int) {
	width := 50.0

	progress := float64(done) / float64(total)
	equals := strings.Repeat("=", int(progress*width))
	dashes := strings.Repeat("-", int((1.0-progress)*width))
	fmt.Print(fmt.Sprintf("["+u.Purple+"*"+u.Normal+"] Finding diffs... [ "+equals+dashes+" ] [%d/%d compared] - (%f", int(done), int(total), math.Round(progress*100.0*100.0)/100.0) + "%) \r")
}
