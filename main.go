package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

type Response struct {
	Answer string `json:"generated_text"`
}

var apiKey = os.Getenv("HF_TOKEN")

const modelEndpoint = "https://api-inference.huggingface.co/models/your_model"

// QueryHuggingFace sends a question to the Hugging Face API and returns the response
func QueryHuggingFace(question string) (string, error) {
	prompt := fmt.Sprintf(`Pretend you are a magic 8 ball. I will give you scenarios, and you will respond in the way a magic 8 ball would, but make it funny and clever. Here is your question: "%s"`, question)
	input := fmt.Sprintf(`{"inputs": "%s"}`, prompt)
	payload := bytes.NewBuffer([]byte(input))

	req, err := http.NewRequest("POST", modelEndpoint, payload)
	if err != nil {
		return "", fmt.Errorf("Error creating request: %s", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("Error making request:", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Error reading response body:", err)
	}

	var responseObject []Response
	err = json.Unmarshal(body, &responseObject)
	if err != nil {
		return "", fmt.Errorf("Error parsing response body:", err)
	}

	answer := responseObject[0].Answer

	return answer, nil
}

func main() {
	app := tview.NewApplication()

	// Input field for the question
	inputField := tview.NewInputField().SetLabel("Ask the Magic 8-Ball: ")
	outputField := tview.NewTextView().SetDynamicColors(true).SetTextAlign(1)

	inputField.SetBorder(true)
	outputField.SetBorder(true)

	// Function to handle when the Enter key is pressed
	inputField.SetDoneFunc(func(key tcell.Key) {
		if key == tcell.KeyEnter {
			question := inputField.GetText()
			// Query the Hugging Face API
			answer, err := QueryHuggingFace(question)
			if err != nil {
				// Handle the error properly in a real application
				fmt.Println("Error querying the API:", err)
				return
			}

			// Display the answer (consider improving UI/UX in a real application)
			fmt.Println("Magic 8-Ball says:", answer)
			outputField.SetText(answer)
			// app.Stop()
		}
	})

	// Set the root layout and run the application
	root := tview.NewFlex().SetDirection(tview.FlexRow).
		AddItem(outputField, 0, 4, false).
		AddItem(inputField, 0, 1, true)

	if err := app.SetRoot(root, true).Run(); err != nil {
		panic(err)
	}
}
