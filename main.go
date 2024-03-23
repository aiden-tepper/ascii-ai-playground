package main

import (
	"time"

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
	debugMode     = false
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

	go func() {
		time.Sleep(10 * time.Millisecond)
		_, _, w, _ := root.GetRect()
		if w < 60 {
			isMobile = true
			// log.Println("Mobile mode detected:", w)
			app.QueueUpdateDraw(func() {
				subPage.SwitchToPage("eightBall")
			})
		}
	}()

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
