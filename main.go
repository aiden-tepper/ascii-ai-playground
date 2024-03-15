package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

var (
	app           *tview.Application
	debugView     *tview.TextView
	outputView    *tview.TextView
	questionView  *tview.TextView
	eightBallView *tview.TextView
)

const (
	debugMode     = false
	modelEndpoint = "https://api-inference.huggingface.co/models/google/gemma-7b-it"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	log.Print("Loaded .env file")
}

func main() {
	app = setupApp()
	root := setupUI()
	if debugMode {
		setupDebugView(root)
	}
	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
