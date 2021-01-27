package util

import (
	"fmt"
	"math"
	"strings"
)

// prints loader
func PrintLoader(done, total int, msg string) {
	width := 50.0
	progress := float64(done) / float64(total)

	if progress >= 1 {
		fmt.Println(Clear + "\r[" + Purple + "*" + Normal + "] " + msg + Bright + "Done ✔" + Normal)
		return
	}

	equals := strings.Repeat("=", int(progress*width))
	dashes := strings.Repeat("-", int((1.0-progress)*width))

	fmt.Print(fmt.Sprintf("["+Purple+"*"+Normal+"] "+msg+" [ "+Blue+equals+">"+Normal+dashes+" ] [%d/%d] - (%f", int(done), int(total), math.Round(progress*100.0*100.0)/100.0) + "%) \r")
}

func PrintSimpleLoader(done, total int, msg string) {
	width := 50.0

	progress := float64(done) / float64(total)
	equals := strings.Repeat("=", int(progress*width))
	dashes := strings.Repeat("-", int((1.0-progress)*width))
	fmt.Print(fmt.Sprintf("["+Purple+"*"+Normal+"] "+msg+" ["+Blue+equals+">"+Normal+dashes+"]"+Yellow+" (%.2f", math.Round(progress*100.0*100.0)/100.0) + "%)" + Normal + "\r")
}
