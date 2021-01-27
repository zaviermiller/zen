package runner

import "github.com/zaviermiller/zen/diff"

// enum type for different outputs
type OutputType int

const (
	PROMPT OutputType = iota
	RESPONSE
)

// output type that can differentiate between prompt and response outputs
type ZenOutput struct {
	Type    OutputType
	Content string
}

// comparison of two outputs returning a diff object
func (o ZenOutput) CompareDiff(o1 ZenOutput) *diff.ZenDiff {
	// display constant
	// om := []string{"PROMPT", "RESPONSE"}

	// check to make sure the output types are the same (worst case if outputTypes dont match)
	// if o.Type != o1.Type {
	// 	return &ZenDiff{errorType: "OUTPUT TYPE", testDiff: om[o1.outputType], correctDiff: om[o.outputType]}
	// }

	// // basic forward checking of output to ensure match
	// // probably need to improve
	// for i, char := range o.content {
	// 	if char != rune(o1.content[i]) {
	// 		return &zenDiff{errorType: om[o1.outputType] + " CONTENT", testDiff: o1.content[i:], correctDiff: o.content[i:]}
	// 	}
	// }

	// lines matched
	return nil
}
