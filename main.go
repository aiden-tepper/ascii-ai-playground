package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gdamore/tcell/v2"
	"github.com/joho/godotenv"
	"github.com/rivo/tview"
)

type Response struct {
	Answer string `json:"generated_text"`
}

var (
	app       *tview.Application
	debugView *tview.TextView
)

// Function to append messages to the debug view
func debugLog(message string) {
	app.QueueUpdateDraw(func() {
		fmt.Fprintf(debugView, "%s\n", message)
	})
}

func showLoadingAnimation(app *tview.Application, outputField *tview.TextView, done chan bool) {
	animationFrames := []string{"Loading .  ", "Loading .. ", "Loading ..."}
	for {
		select {
		case <-done:
			return
		default:
			for _, frame := range animationFrames {
				app.QueueUpdateDraw(func() {
					outputField.SetText(frame)
				})
				time.Sleep(200 * time.Millisecond)
			}
		}
	}
}

const modelEndpoint = "https://api-inference.huggingface.co/models/google/gemma-7b-it"

// QueryHuggingFace sends a question to the Hugging Face API and returns the response
func QueryHuggingFace(question string) (string, error) {
	apiKey := os.Getenv("HF_TOKEN")
	prompt := fmt.Sprintf(`Pretend you are a magic 8 ball. I will give you scenarios, and you will respond in the way a magic 8 ball would, but make it funny and clever. Here is your question: "%s"`, question)
	input := fmt.Sprintf(`{"inputs": "%s"}`, prompt)
	payload := bytes.NewBuffer([]byte(input))

	req, err := http.NewRequest("POST", modelEndpoint, payload)
	if err != nil {
		return "", fmt.Errorf("error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %s", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %s", err)
	}

	var responseObject []Response
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return "", fmt.Errorf("error parsing response body: %s", err)
	}

	answer := responseObject[0].Answer

	return answer, nil
}

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
	log.Print("Loaded .env file")
}

func main() {
	app = tview.NewApplication()

	// Input field for the question
	inputField := tview.NewInputField().SetLabel("Ask the Magic 8-Ball: ")
	outputField := tview.NewTextView().SetDynamicColors(true).SetTextAlign(1)

	inputField.SetBorder(true)
	outputField.SetBorder(true)

	// Function to handle when the Enter key is pressed
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			question := inputField.GetText()
			done := make(chan bool)

			go showLoadingAnimation(app, outputField, done)

			// Use a goroutine for querying the API to not block the main thread
			go func() {
				answer, err := QueryHuggingFace(question)
				if err != nil {
					// Proper error handling
					fmt.Println("Error querying the API:", err)
					done <- true
					return
				}

				app.QueueUpdateDraw(func() {
					outputField.SetText("Magic 8-Ball says: " + answer)
				})
				done <- true
			}()
		}
	})

	debugView = tview.NewTextView()
	debugView.SetTitle("Debug Log").SetBorder(true)

	// Set the root layout and run the application
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(outputField, 0, 4, false).
		AddItem(inputField, 0, 1, true).
		AddItem(debugView, 0, 1, false)

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
