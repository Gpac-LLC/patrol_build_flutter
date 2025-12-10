package print

import (
	"fmt"
	"runtime"
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
)

func init() {
	if runtime.GOOS == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
	}
}

func _printColor(colorCode string, message string) {
	fmt.Println(colorCode + message + Reset)
}

func Error(message string) {
	_printColor(Red, message)
}

func Success(message string) {
	_printColor(Green, message)
}

func Warning(message string) {
	_printColor(Yellow, message)
}

func Action(message string) {
	_printColor(Blue, message)
}

func StepCompleted(message string) {
	_printColor(Purple, message)
}

func StepInitiated(message string) {
	_printColor(Cyan, message)
}

func Vanilla(message string) {
	fmt.Println(message)
}
