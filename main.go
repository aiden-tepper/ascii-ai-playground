package main

import (
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

var (
	app           *tview.Application
	outputView    *tview.TextView
	questionView  *tview.TextView
	eightBallView *tview.TextView
)

const (
	debugMode     = true
	modelEndpoint = "https://api-inference.huggingface.co/models/google/gemma-7b-it"
)

func init() {
	if err := godotenv.Load(); err != nil {
		alert("No .env file found")
	}
}

func main() {
	app = setupApp()
	root := setupUI()
	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
