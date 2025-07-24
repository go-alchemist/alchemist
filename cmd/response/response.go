package response

import (
	"os"

	"github.com/fatih/color"
)

func Error(err error) {
	if err != nil {
		color.Red(err.Error())
		os.Exit(1)
	}
}

func Success(message string) {
	if message != "" {
		color.Green(message)
	}
	os.Exit(0)
}

func Info(message string) {
	if message != "" {
		color.Blue(message)
	}
}
func Warning(message string) {
	if message != "" {
		color.Yellow(message)
	}
}
