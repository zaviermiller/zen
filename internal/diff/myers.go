// MYERS DIFF ALGORITHM IMPLEMENTED BY FOLLOWING THIS BLOG POST
// https://blog.jcoglan.com/2017/02/17/the-myers-diff-algorithm-part-3/

package diff

import (
	"fmt"
	"strings"

	d "github.com/zaviermiller/zen/internal/display"
)

type MOperation int

const (
	EQL MOperation = iota
	INS
	DEL
)

type MyersDiff struct {
	CorrectName    string
	TestName       string
	CorrectContent []string
	TestContent    []string
	Ops            []MOperation
}

// Calculate figures out the meyers diff traces for the given string slices
func (md *MyersDiff) Calculate(correct, test []string) float64 {

	max := max(len(correct), len(test))

	correctArr := []string{}
	testArr := []string{}

	for i := 0; i < max; i++ {
		if i < len(correct) {
			correctArr = append(correctArr, strings.Split(correct[i], "")...)
		}
		if i < len(test) {
			testArr = append(testArr, strings.Split(test[i], "")...)
		}
	}

	md.CorrectContent = correctArr
	md.TestContent = testArr

	var score float64
	md.Ops, score = shortestEditScript(testArr, correctArr)

	return score

}

// Print goes thru and prints out the meyers diffs
func (md MyersDiff) Print() {
	srcIndex, dstIndex := 0, 0
	score := 0

	fmt.Print(d.Clear)
	printStr := ""

	for _, op := range md.Ops {
		switch op {
		case INS:
			printStr += (d.GreenHighlight + d.White + d.Bright + (md.CorrectContent[dstIndex]) + d.Normal)
			dstIndex += 1
			score++

		case EQL:
			printStr += md.CorrectContent[dstIndex]
			srcIndex += 1
			dstIndex += 1

		case DEL:
			printStr += d.RedHighlight + d.White + d.Bright + (md.TestContent[srcIndex]) + d.Normal
			srcIndex += 1
			score++
		}
	}

	// if printStr != printStr {
	// 	fmt.Println(d.Red + d.Bright + d.Clear + "!!!DIFF DETECTED!!! " + d.Normal + "\n")
	// 	fmt.Println("   [" + d.Blue + md.CorrectName + d.Normal + "]: " + d.Green + printStr + d.Normal)
	// 	fmt.Println("   [" + d.Blue + md.TestName + d.Normal + "]: " + d.Red + printStr + d.Normal + "\n")
	// }

	fmt.Println(printStr)
}

// func (md MyersDiff) shortestEdit() [][]int {
// 	n, m := len(md.CorrectContent), len(md.TestContent)
// 	if n > m {
// 		n, m = m, n
// 	}

// 	max := n + m

// 	// v[1] = 0
// 	trace := make([][]int, 2*max+1)

// 	for d := 0; d < max; d++ {
// 		v := make([]int, 2*max+1)
// 		trace = append(trace, v)
// 		// if d == 0 {
// 		// 	subLen := 0
// 		// 	for len(md.CorrectContent) > subLen && len(md.TestContent) > subLen && md.CorrectContent == md.TestContent {
// 		// 		subLen++
// 		// 	}
// 		// 	if subLen == len(md.CorrectContent) && subLen == len(md.TestContent) {
// 		// 		// same content
// 		// 		return nil
// 		// 	}
// 		// }
// 		dOffset := max

// 		for k := -d; k <= d; k += 2 {
// 			vCopy := make([]int, len(v))

// 			var x int
// 			if k == -d || (k != d && v[dOffset+k-1] < v[dOffset+k+1]) {
// 				x = v[k+1]
// 			} else {
// 				x = v[k-1]
// 			}

// 			y := x - k

// 			for x < n && y < m && md.CorrectContent[x] == md.TestContent[x] {
// 				x, y = x+1, y+1
// 			}

// 			v[k + dOffset] = x

// 			if x >= n && y >= m {
// 				copy(vCopy, v)
// 				trace[d] = vCopy
// 				return trace, dOffset
// 			}
// 		}
// 	}
// 	return nil, 0
// }

// func (md MyersDiff) backtrack(trace [][]int, offset int) {
// 	x, y := len(md.CorrectContent), len(md.TestContent)
// 	paths := make([][]int, len(trace))

// 	d := len(trace) - 1

// 	var prevType MOperation
// 	var substr string

// 	for x > 0 && y > 0 && d > 0; d-- {
// 		v := trace[d]
// 		if len(v) == 0 {
// 			// same i think
// 			if prevType == EQL {
// 				substr +=
// 			} else {

// 			}
// 			continue
// 		}

// 		k := x - y

// 		var prevK int
// 		if k == -d || (k != d && v[offset+k-1] < v[offset+k+1]) {
// 			prevK = k + 1
// 		} else {
// 			prevK = k - 1
// 		}

// 		prevX := v[prevK + offset]
// 		prevY := prevX - prevK

// 		for x > prevX && y > prevY {
// 			dChan <- MyersData{x: x, y: y, prevX: prevX, prevY: prevY}
// 			x, y = x-1, y-1
// 		}

// 		if d > 0 {
// 			dChan <- MyersData{x: x, y: y, prevX: prevX, prevY: prevY}
// 		}

// 		x, y = prevX, prevY

// 	}

// }

// func reverse(s [][]int) [][]int {
// 	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
// 		s[i], s[j] = s[j], s[i]
// 	}
// 	return s
// }

// TEMP RIP FROM https://github.com/cj1128/myers-diff/blob/master/main.go

func shortestEditScript(src, dst []string) ([]MOperation, float64) {
	n := len(src)
	m := len(dst)
	max := n + m
	var trace []map[int]int
	var x, y int

loop:
	for d := 0; d <= max; d++ {
		// 最多只有 d+1 个 k
		v := make(map[int]int, d+2)
		trace = append(trace, v)

		// 需要注意处理对角线
		if d == 0 {
			t := 0
			for len(src) > t && len(dst) > t && src[t] == dst[t] {
				t++
			}
			v[0] = t
			if t == len(src) && t == len(dst) {
				break loop
			}
			continue
		}

		lastV := trace[d-1]

		for k := -d; k <= d; k += 2 {
			// 向下
			if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
				x = lastV[k+1]
			} else { // 向右
				x = lastV[k-1] + 1
			}

			y = x - k

			// 处理对角线
			for x < n && y < m && src[x] == dst[y] {
				x, y = x+1, y+1
			}

			v[k] = x

			if x == n && y == m {
				break loop
			}
		}
	}

	// for debug
	// printTrace(trace)

	// 反向回溯
	var script []MOperation

	diffCount := 0
	x = n
	y = m
	var k, prevK, prevX, prevY int

	for d := len(trace) - 1; d > 0; d-- {
		k = x - y
		lastV := trace[d-1]

		if k == -d || (k != d && lastV[k-1] < lastV[k+1]) {
			prevK = k + 1
		} else {
			prevK = k - 1
		}

		prevX = lastV[prevK]
		prevY = prevX - prevK

		for x > prevX && y > prevY {
			script = append(script, EQL)
			x -= 1
			y -= 1
			if diffCount > 0 {
				diffCount--
			}
		}

		if x == prevX {
			script = append(script, INS)
			diffCount++
		} else {
			script = append(script, DEL)
			diffCount++
		}

		x, y = prevX, prevY
	}

	if trace[0][0] != 0 {
		for i := 0; i < trace[0][0]; i++ {
			script = append(script, EQL)
			if diffCount > 0 {
				diffCount--
			}
		}
	}

	return reverse(script), 1 - (float64(diffCount) / float64(len(script)))
}

func reverse(s []MOperation) []MOperation {
	result := make([]MOperation, len(s))

	for i, v := range s {
		result[len(s)-1-i] = v
	}

	return result
}
