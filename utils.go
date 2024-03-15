package main

import (
	"fmt"
	"strings"
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

func debugLog(message string) {
	switch debugMode {
	case true:
		app.QueueUpdateDraw(func() {
			fmt.Fprintf(debugView, "%s\n", message)
		})
	case false:
		fmt.Fprintf(debugView, "%s\n", message)
	}
}
