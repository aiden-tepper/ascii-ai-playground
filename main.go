package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type Response struct {
	Answer string `json:"generated_text"`
}

var (
	app        *tview.Application
	debugView  *tview.TextView
	outputView *tview.TextView
)

const (
	debugMode      = false
	modelEndpoint  = "https://api-inference.huggingface.co/models/google/gemma-7b-it"
	eightBallAscii = `	   ____
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
)

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

// Append messages to the debug view
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

// Animate a loading message: "Loading .  ", "Loading .. ", "Loading ...", repeat
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

// QueryHuggingFace sends a question to the Hugging Face API and returns the response
func QueryHuggingFace(question string) (map[string]string, error) {
	apiKey := os.Getenv("HF_TOKEN")
	prompt := fmt.Sprintf(`Pretend you are a magic 8 ball. I will give you a question, and you will respond in the way a magic 8 ball would, but make it funny and clever. Here is your question: '%s'. Reply in this format: {\"answer\": answer, \"explanation\": explanation}, where 'answer' is the few word answer that would show up on the magic 8-ball itself, and explanation is a sentence or two of explanation, humorous quips, or highly analytical statements.`, question)
	input := fmt.Sprintf(`{"inputs": "%s"}`, prompt)
	payload := bytes.NewBuffer([]byte(input))

	req, err := http.NewRequest("POST", modelEndpoint, payload)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %s", err)
	}

	var responseObject []Response
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return nil, fmt.Errorf("error parsing response body: %s", err)
	}

	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(responseObject[0].Answer, -1)

	var result map[string]string
	err = json.Unmarshal([]byte(matches[1][0]), &result)
	if err != nil {
		return nil, fmt.Errorf("error parsing JSON object: %s", err)
	}

	if len(responseObject) > 0 {
		return map[string]string{"answer": result["answer"], "explanation": result["explanation"]}, nil
	}
	return nil, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	log.Print("Loaded .env file")
}

func main() {
	app = tview.NewApplication()

	inputField := tview.NewInputField().SetLabel("Ask the Magic 8-Ball: ")
	outputView = tview.NewTextView().SetDynamicColors(true).SetTextAlign(1)
	questionView := tview.NewTextView().SetTextAlign(1).SetTextStyle(tcell.StyleDefault.Italic(true))
	eightBallView := tview.NewTextView().SetTextAlign(0).SetText(eightBallAscii)

	xOffset, yOffset := analyzeMultilineString(eightBallAscii)
	log.Printf("xOffset: %d, yOffset: %d", xOffset/2, yOffset/2)

	eightBallView.SetDrawFunc(func(screen tcell.Screen, x, y, w, h int) (int, int, int, int) {
		y += h / 2
		x += w / 2
		xOffset, yOffset := analyzeMultilineString(eightBallAscii)
		return x - xOffset/2, y - yOffset/2, w, h
	})

	inputField.SetBorder(true)
	outputView.SetBorder(true)
	questionView.SetBorder(true)
	eightBallView.SetBorder(true)

	done := make(chan bool)

	// Function to handle when the Enter key is pressed
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			question := inputField.GetText()
			inputField.SetText("")
			questionView.SetText(question)
			showLoadingAnimation(done)

			// Use a goroutine for querying the API to not block the main thread
			go func() {
				defer func() {
					done <- true
				}()

				result, err := QueryHuggingFace(question)
				if err != nil {
					if debugMode {
						debugLog(fmt.Sprintf("Error querying the API: %s", err))
					} else {
						log.Printf("Error querying the API: %s", err)
					}
					return
				}

				app.QueueUpdateDraw(func() {
					output := fmt.Sprint("[::b]" + result["answer"] + "\n\n[::Bi]" + result["explanation"])
					_, yOffset := analyzeMultilineString(output)
					_, _, _, height := outputView.GetInnerRect()
					outputView.SetText(strings.Repeat("\n", height/2-yOffset/2) + output)
				})
			}()
		}
	})

	subView := tview.NewFlex().SetDirection(tview.FlexColumn).AddItem(eightBallView, 0, 1, false).AddItem(outputView, 0, 1, false)

	contentView := tview.NewFlex().SetDirection(tview.FlexRow).AddItem(questionView, 3, 0, false).AddItem(subView, 0, 4, false)
	contentView.SetBorder(true)

	// Set the root layout and run the application
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(contentView, 0, 4, false).
		AddItem(inputField, 3, 0, true)

	if debugMode {
		debugView = tview.NewTextView()
		debugView.SetTitle("Debug Log").SetBorder(true)
		root.AddItem(debugView, 0, 1, false)
	}

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
