package main

import (
	"log"
	"strings"

	"github.com/rivo/tview"
)

type Response struct {
	Answer string `json:"generated_text"`
}

func analyzeMultilineString(s string) (maxLength int, lineCount int) {
	lines := strings.Split(s, "\n") // Split the string into lines
	lineCount = len(lines)          // The number of lines is the length of the slice

	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line) // Update maxLength if the current line is longer
		}
	}

	return maxLength, lineCount
}

func alert(message string) {
	switch debugMode {
	case true:
		log.Println(message)
		app.Stop()
	case false:
		modal := tview.NewModal().
			SetText(message).
			AddButtons([]string{"Restart", "Quit"}).
			SetDoneFunc(func(buttonIndex int, buttonLabel string) {
				if buttonLabel == "Quit" {
					app.Stop()
				} else if buttonLabel == "Restart" {
					main()
				}
			})
		if err := app.SetRoot(modal, false).SetFocus(modal).Run(); err != nil {
			panic(err)
		}

	}
}
