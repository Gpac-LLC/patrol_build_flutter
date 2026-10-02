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

func SetColorsForOS(goos string) {
	if goos == "windows" {
		Reset = ""
		Red = ""
		Green = ""
		Yellow = ""
		Blue = ""
		Purple = ""
		Cyan = ""
	}
}

func init() {
	SetColorsForOS(runtime.GOOS)
}

func printColor(colorCode string, message string) {
	fmt.Println(colorCode + message + Reset)
}

func Error(message string) {
	printColor(Red, message)
}

func Errorf(format string, a ...any) {
	printColor(Red, fmt.Sprintf(format, a...))
}

func Success(message string) {
	printColor(Green, message)
}

func Successf(format string, a ...any) {
	printColor(Green, fmt.Sprintf(format, a...))
}

func Warning(message string) {
	printColor(Yellow, message)
}

func Warningf(format string, a ...any) {
	printColor(Yellow, fmt.Sprintf(format, a...))
}

func Info(message string) {
	printColor(Cyan, message)
}

func Infof(format string, a ...any) {
	printColor(Cyan, fmt.Sprintf(format, a...))
}

func Action(message string) {
	printColor(Blue, message)
}

func Actionf(format string, a ...any) {
	printColor(Blue, fmt.Sprintf(format, a...))
}

func StepCompleted(message string) {
	printColor(Purple, message)
}

func StepCompletedf(format string, a ...any) {
	printColor(Purple, fmt.Sprintf(format, a...))
}

func StepInitiated(message string) {
	printColor(Cyan, message)
}

func StepInitiatedf(format string, a ...any) {
	printColor(Cyan, fmt.Sprintf(format, a...))
}

func Vanilla(message string) {
	fmt.Println(message)
}

func Vanillaf(format string, a ...any) {
	fmt.Println(fmt.Sprintf(format, a...))
}
