package util

import (
	"fmt"
)

var Red = "\033[31m"

var Green = "\033[32m"

var Yellow = "\033[33m"

var Blue = "\033[34m"

var Purple = "\033[35m"

var Bright = "\033[1m"

var Dark = "\033[2m"

var Normal = "\033[0m"

var Clear = "\u001b[2K"

func Zen(v string) string {
	return ` ____         
/_  / ___ ___ 
 / /_/ -_) _ \
/___/\__/_//_/ v` + v
}

func PrintLogo(v string) {
	// pretty stuff
	fmt.Println(Green + Bright + Zen(v))
	fmt.Println(Purple + "\n   \"Bring yourself peace\" -Z" + Normal + "\n")
}

// log helpers
func ZenLog(msg string) {
	fmt.Print("\n\r[" + Purple + "*" + Normal + "] " + msg)
}

func ZenLogLn(msg string) {
	fmt.Println("\n\r[" + Purple + "*" + Normal + "] " + msg)
}

func ZenWeirdLog(msg string) {
	fmt.Print("\r[" + Purple + "*" + Normal + "] " + msg)
}

// computes color of score
func ColorScore(score float64) string {
	textScore := fmt.Sprintf("%.2f", score)
	if score >= .75 {
		return Green + textScore + Normal
	} else if score >= .5 {
		return Yellow + textScore + Normal
	} else {
		return Red + textScore + Normal
	}
}
