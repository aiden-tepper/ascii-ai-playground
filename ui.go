package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

var shakeDegreeX, shakeDegreeY = 0, 0

const eightBallAscii = `	   ____
    ,dP9CGG88@b,
  ,IP  _   Y888@@b,
 dIi  (_)   G8888@b
dCII  (_)   G8888@@b
GCCIi     ,GG8888@@@
GGCCCCCCCGGG88888@@@
GGGGCCCGGGG88888@@@@...
Y8GGGGGG8888888@@@@P.....
 Y88888888888@@@@@P......
 'Y8888888@@@@@@@P'......
    '@@@@@@@@@P'.......
        """"........`

func setupApp() *tview.Application {
	return tview.NewApplication()
}

func setupUI() *tview.Flex {
	inputField := tview.NewInputField().SetLabel("Ask the Magic 8-Ball: ")
	outputView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(1)
	questionView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(1)
	eightBallView = tview.NewTextView().SetTextAlign(0).SetText(eightBallAscii)

	inputField.SetBorder(true).SetBorderPadding(0, 0, 2, 2)
	outputView.SetBorder(true).SetBorderPadding(4, 4, 4, 4)
	questionView.SetBorder(true).SetBorderPadding(0, 0, 2, 2)
	eightBallView.SetBorder(true).SetBorderPadding(4, 4, 4, 4)

	inputField.SetDoneFunc(func(key tcell.Key) {
		handleInput(key, inputField)
	})

	eightBallView.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h/2 + shakeDegreeY
		x += w/2 + shakeDegreeX
		xOffset, yOffset := analyzeMultilineString(eightBallAscii)
		return x - xOffset/2, y - yOffset/2, w, h
	})

	subView := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView().SetDynamicColors(true), 3, 0, false).
		AddItem(eightBallView, 0, 1, false).
		AddItem(tview.NewTextView().SetDynamicColors(true), 3, 0, false).
		AddItem(outputView, 0, 1, false).
		AddItem(tview.NewTextView().SetDynamicColors(true), 3, 0, false)

	questionViewBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(tview.NewTextView().SetDynamicColors(true), 3, 0, false).
		AddItem(questionView, 0, 1, false).
		AddItem(tview.NewTextView().SetDynamicColors(true), 3, 0, false)

	contentView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(tview.NewTextView().SetDynamicColors(true), 1, 0, false).
		AddItem(questionViewBox, 3, 0, false).
		AddItem(tview.NewTextView().SetDynamicColors(true), 1, 0, false).
		AddItem(subView, 0, 4, false).
		AddItem(tview.NewTextView().SetDynamicColors(true), 1, 0, false)
	contentView.SetBorder(true)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(contentView, 0, 4, false).
		AddItem(inputField, 3, 0, true)
}

func handleInput(key tcell.Key, inputField *tview.InputField) {
	if key == tcell.KeyEnter {
		doneLoading := make(chan bool)
		doneShaking := make(chan bool)
		question := inputField.GetText()
		inputField.SetText("")
		questionView.SetText("[::i]" + question)
		showLoadingAnimation(doneLoading)
		showShakingAnimation(doneShaking)
		queryAPIAndDisplayResult(question, doneLoading, doneShaking)
	}
}

func showLoadingAnimation(done chan bool) {
	go func() {
		animationFrames := []string{"Loading .  ", "Loading .. ", "Loading ..."}
		frameIndex := 0
		for {
			select {
			case <-done:
				return
			default:
				app.QueueUpdateDraw(func() {
					_, _, _, height := outputView.GetInnerRect()
					outputView.SetText(strings.Repeat("\n", height/2) + animationFrames[frameIndex])
				})
				frameIndex = (frameIndex + 1) % len(animationFrames)
				time.Sleep(200 * time.Millisecond)
			}
		}
	}()
}

func showShakingAnimation(done chan bool) {
	go func() {
		for {
			select {
			case <-done:
				shakeDegreeX, shakeDegreeY = 0, 0
				app.QueueUpdateDraw(func() {
					eightBallView.SetText(eightBallAscii)
				})
				return
			default:
				app.QueueUpdateDraw(func() {
					shakeDegreeX = rand.Intn(4) - 2
					shakeDegreeY = rand.Intn(4) - 2
					eightBallView.SetText(eightBallAscii)
				})
				time.Sleep(50 * time.Millisecond)
			}
		}
	}()
}
