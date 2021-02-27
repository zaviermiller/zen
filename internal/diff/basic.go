package diff

import (
	"fmt"
	"strconv"

	d "github.com/zaviermiller/zen/internal/display"
)

// BasicDiff is a difference object that contains data for printing
// the diffs to the screen -- very basic, and what v1.0 used :]
// THIS IS DEPRECATED PLEASE DO NOT USE IT. it will be removed soon.
type BasicDiff struct {
	CorrectID   int
	TestID      int
	CorrectName string
	TestName    string

	// 1 2 3ab d 5 6 <-- calc 2 worksheet nums
	TestDiff    string
	CorrectDiff string
	Type        string
}

// func NewBasicDiff(cid, tid int, cname, tname string) BasicDiff {
// 	diff := BasicDiff{CorrectID: cid, TestID: tid, CorrectName: cname, TestName: tname}

// 	return diff
// }

// func (diff BasicDiff) Calculate(o1, o2 string) error {
// 	// display constant
// 	// om := []string{"PROMPT", "RESPONSE"}

// 	// basic forward checking of output to ensure match
// 	for i, char := range o1 {
// 		if char != rune(o2[i]) {
// 			diff.TestDiff = o2[i:]
// 			diff.CorrectDiff = o1[:i]
// 			return nil
// 		}
// 	}

// 	// lines matched
// 	return ZenNoDiffError{}
// }

// stay awake and here until the drug wears off its the worst thing and thinkng about it makes it happen <- this is the truth and while they may not know it now its what everything returns to. idk
// so stop asking the question lol just keep coding??? im not trying to learning anything is no lesson idk

// dont follow the "whatever it is..." ~~ AND STOP TRYING TO THINK ABOUT WHATEVER

// keep trying to describe "it" but noone will know lol. its not forever

// HEY YOU BIG MAN- QUIT - lol im kindaaaa going crazy but

// type some more stuff the trip is over // always remember that ok? thats what makes sense

// wonder why this is the one? theres no reason it is literally just the drug lol
// thats so funny me, keep on moving i guess

// Calculate is DEPRECATED, USE MeyersDiff Type
func (diff BasicDiff) Calculate(out1, out2 []string) float64 {
	// for i, char := range diff. {
	// 	if char != rune(o2[i]) {
	// 		diff.TestDiff = o2[i:]
	// 		diff.CorrectDiff = o1[:i]
	// 		return nil
	// 	}
	// }
	return 0
}

// Print out the diff
func (z BasicDiff) Print() {
	fmt.Println(d.Red + d.Bright + d.Clear + "!!!" + z.Type + " DIFF!!! " + d.Normal + "\n")
	fmt.Println("   " + z.CorrectName + " [" + strconv.Itoa(z.CorrectID) + "]: " + d.Green + z.CorrectDiff + d.Normal)
	fmt.Println("   " + z.TestDiff + " [" + strconv.Itoa(z.TestID) + "]: " + d.Red + z.TestDiff + d.Normal + "\n")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
