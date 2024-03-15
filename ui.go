package main

import (
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

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
	questionView = tview.NewTextView().SetTextAlign(1).SetTextStyle(tcell.StyleDefault.Italic(true))
	eightBallView = tview.NewTextView().SetTextAlign(0).SetText(eightBallAscii)

	inputField.SetBorder(true).SetBorderPadding(0, 0, 2, 2)
	outputView.SetBorder(true).SetBorderPadding(4, 4, 4, 4)
	questionView.SetBorder(true).SetBorderPadding(0, 0, 2, 2)
	eightBallView.SetBorder(true).SetBorderPadding(4, 4, 4, 4)

	done := make(chan bool)
	inputField.SetDoneFunc(func(key tcell.Key) {
		handleInput(key, inputField, done)
	})

	eightBallView.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h / 2
		x += w / 2
		xOffset, yOffset := analyzeMultilineString(eightBallAscii)
		return x - xOffset/2, y - yOffset/2, w, h
	})

	subView := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 3, 0, false).
		AddItem(eightBallView, 0, 1, false).
		AddItem(nil, 3, 0, false).
		AddItem(outputView, 0, 1, false).
		AddItem(nil, 3, 0, false)

	questionViewBox := tview.NewFlex().SetDirection(tview.FlexColumn).
		AddItem(nil, 3, 0, false).
		AddItem(questionView, 0, 1, false).
		AddItem(nil, 3, 0, false)

	contentView := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(nil, 1, 0, false).
		AddItem(questionViewBox, 3, 0, false).
		AddItem(nil, 1, 0, false).
		AddItem(subView, 0, 4, false).
		AddItem(nil, 1, 0, false)
	contentView.SetBorder(true)

	return tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(contentView, 0, 4, false).
		AddItem(inputField, 3, 0, true)
}

func handleInput(key tcell.Key, inputField *tview.InputField, done chan bool) {
	if key == tcell.KeyEnter {
		question := inputField.GetText()
		inputField.SetText("")
		questionView.SetText(question)
		showLoadingAnimation(done)
		queryAPIAndDisplayResult(question, done)
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
