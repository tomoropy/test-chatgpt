package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/r3labs/sse/v2"
)

type Choice struct {
	Delta struct {
		Content string `json:"content"`
	} `json:"delta"`
	Index        int         `json:"index"`
	FinishReason interface{} `json:"finish_reason"`
}

type JSONData struct {
	ID      string   `json:"id"`
	Object  string   `json:"object"`
	Created int64    `json:"created"`
	Model   string   `json:"model"`
	Choices []Choice `json:"choices"`
}

const APIEndpoint = "https://api.openai.com/v1/chat/completions"

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type requestBody struct {
	Model    string    `json:"model"`
	Stream   bool      `json:"stream"`
	Messages []Message `json:"messages"`
}

type customTransport struct {
	http.RoundTripper
}

func main() {
	client := &http.Client{
		Transport: &customTransport{
			RoundTripper: http.DefaultTransport,
		},
	}

	messages := []Message{
		{Role: "user", Content: "Server-Sent Eventについて教えて"},
	}

	body := requestBody{
		Messages: messages,
		Model:    "gpt-3.5-turbo",
		Stream:   true,
	}
	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshaling request body:", err)
		os.Exit(1)
	}

	sseClient := sse.NewClient(APIEndpoint)
	sseClient.Connection = client
	sseClient.Method = "POST"
	sseClient.Body = bytes.NewBuffer([]byte(jsonData))
	sseClient.SubscribeRaw(func(msg *sse.Event) {
		var jsonData JSONData
		err := json.Unmarshal([]byte(msg.Data), &jsonData)
		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Printf("%s", jsonData.Choices[0].Delta.Content)
	})

}

func (t *customTransport) RoundTrip(req *http.Request) (*http.Response, error) {

	API_KEY := os.Getenv("API_KEY")
	if API_KEY == "" {
		fmt.Println("API_KEY is not set")
		os.Exit(1)
	}

	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", API_KEY))
	req.Header.Set("Content-Type", "application/json")

	resp, err := t.RoundTripper.RoundTrip(req)
	if err != nil {
		return nil, err
	}

	return resp, err
}
