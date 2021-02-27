package diff

import (
	"math"
	"testing"
)

// test the myers diff algorithm
func TestMyersDiffCalculate(t *testing.T) {

	correct := []string{"zavier miller is dope."}
	test := []string{"Zavier Miller is cool!"}

	myers := NewDiff(MYERS, "Correct", "Test")

	score := math.Round(myers.Calculate(correct, test) * 100)

	if score != 93 {
		t.Errorf("MyersDiff Calculate([]string{\"zavier miller is dope.\"}, []string{\"Zavier Miller is cool!\"}) FAILED, expected score %d, but got %v\n\n", 93, score)
	} else {
		t.Logf("MyersDiff Calculate([]string{\"zavier miller is dope.\"}, []string{\"Zavier Miller is cool!\"}) PASSED, expected %d and got %v\n\n", 93, score)
		myers.Print()
	}

	perfScore := math.Round(myers.Calculate(correct, correct) * 100)

	if perfScore != 100 {
		t.Errorf("MyersDiff Calculate([]string{\"zavier miller is dope.\"}, []string{\"Zavier Miller is cool!\"}) FAILED, expected score %d, but got %v\n\n", 100, perfScore)
	} else {
		t.Logf("MyersDiff Calculate([]string{\"zavier miller is dope.\"}, []string{\"zavier miller is dope.\"}) PASSED, expected %d and got %v\n\n", 100, perfScore)
	}

}
