package diff

import (
	"fmt"
	"strconv"

	u "github.com/zaviermiller/zen/internal/util"
)

// difference object that contains data for printing
// the diffs to the screen
type ZenDiff struct {
	TestDiff    string
	CorrectDiff string
	Type        string
}

// prints out the diff
func (z ZenDiff) Print(correctId, testId int, correctBinPath, testBinPath string) {
	fmt.Println(u.Red + u.Bright + u.Clear + "!!!" + z.Type + " DIFF!!! " + u.Normal + "\n")
	fmt.Println("   " + correctBinPath + " [" + strconv.Itoa(correctId) + "]: " + u.Green + z.CorrectDiff + u.Normal)
	fmt.Println("   " + testBinPath + " [" + strconv.Itoa(testId) + "]: " + u.Red + z.TestDiff + u.Normal + "\n")
}
