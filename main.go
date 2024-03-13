package main

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// Replace with your actual Hugging Face API key
const apiKey = "your_hugging_face_api_key"

// Replace with the chosen LLM model endpoint
const modelEndpoint = "https://api-inference.huggingface.co/models/your_model"

// QueryHuggingFace sends a question to the Hugging Face API and returns the response
func QueryHuggingFace(question string) (string, error) {
    // Prepare the request body with the question
    // requestBody, err := json.Marshal(map[string]string{
    //     "inputs": question,
    // })
    // if err != nil {
    //     return "", err
    // }

    // // Create a new HTTP request
    // req, err := http.NewRequest("POST", modelEndpoint, bytes.NewBuffer(requestBody))
    // if err != nil {
    //     return "", err
    // }

    // // Set the required headers
    // req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", apiKey))
    // req.Header.Set("Content-Type", "application/json")

    // // Make the HTTP request
    // client := &http.Client{}
    // resp, err := client.Do(req)
    // if err != nil || resp.StatusCode != 200 {
    //     return "", err
    // }
    // defer resp.Body.Close()

    // // Read and parse the response body
    // responseBody, err := ioutil.ReadAll(resp.Body)
    // if err != nil {
    //     return "", err
    // }

    // // Extract the answer (adjust according to the API's response structure)
    // var responseMap map[string]interface{}
    // if err := json.Unmarshal(responseBody, &responseMap); err != nil {
    //     return "", err
    // }

    // // Assuming the answer is in the text field
    // answer, ok := responseMap["generated_text"].(string)
    // if !ok {
    //     return "", fmt.Errorf("invalid response format")
    // }

    // return answer, nil
    return question, nil
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
