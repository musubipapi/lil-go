package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
)

// Message represents a single message in the conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatCompletionRequest represents the request structure for the chat completion API
type ChatCompletionRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
	Stream   *bool     `json:"stream,omitempty"`
}

type ChatCompletionStreamResponse struct {
	ID                string `json:"id"`
	Object            string `json:"object"`
	Created           int64  `json:"created"`
	Model             string `json:"model"`
	SystemFingerprint string `json:"system_fingerprint"`
	Choices           []struct {
		Index int `json:"index"`
		Delta struct {
			Role    string `json:"role,omitempty"`
			Content string `json:"content,omitempty"`
		} `json:"delta"`
		FinishReason *string `json:"finish_reason"`
	} `json:"choices"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	var rootCmd = &cobra.Command{Use: "cl(a)i"}
	var gptCmd = &cobra.Command{
		Use:   "gpt [prompt]",
		Short: "Generate a response using GPT-4o",
		Args:  cobra.ArbitraryArgs,
		Run: func(cmd *cobra.Command, args []string) {
			prompt := strings.Join(args, " ")
			if prompt == "" {
				fmt.Println("Error: Please provide a prompt.")
				fmt.Println("Usage: cl(a)i gpt [prompt]")
				os.Exit(1)
			}
			apiKey := os.Getenv("OPENAI_KEY")
			apiUrl := "https://api.openai.com/v1/chat/completions"
			stream := true
			data := ChatCompletionRequest{
				Model: "gpt-4o-mini",
				Messages: []Message{
					{
						Role:    "system",
						Content: "You are a helpful assistant.",
					},
					{
						Role:    "user",
						Content: prompt,
					},
				},
				Stream: &stream,
			}
			jsonData, err := json.Marshal(data)
			if err != nil {
				fmt.Println("Error encoding JSON:", err)
				return
			}

			req, err := http.NewRequest("POST", apiUrl, bytes.NewBuffer(jsonData))
			if err != nil {
				fmt.Println("Error creating request:", err)
				return
			}
			// Set headers
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+apiKey)

			// Send the request
			client := &http.Client{}
			resp, err := client.Do(req)
			if err != nil {
				fmt.Println("Error sending request:", err)
				return
			}
			defer resp.Body.Close()

			// Check for HTTP errors
			if resp.StatusCode != http.StatusOK {
				fmt.Printf("HTTP request failed with status %d\n", resp.StatusCode)
				return
			}
			reader := bufio.NewReader(resp.Body)
			fmt.Println("Streaming response:")
			for {
				line, err := reader.ReadBytes('\n')
				if err != nil {
					if err == io.EOF {
						break // End of stream
					}
					fmt.Println("Error reading stream:", err)
					return
				}

				line = bytes.TrimSpace(line)
				if !bytes.HasPrefix(line, []byte("data: ")) {
					continue
				}

				data := bytes.TrimPrefix(line, []byte("data: "))
				if string(data) == "[DONE]" {
					break
				}

				var streamResp ChatCompletionStreamResponse
				if err := json.Unmarshal(data, &streamResp); err != nil {
					fmt.Println("Error unmarshalling JSON:", err)
					continue
				}

				if len(streamResp.Choices) > 0 {
					content := streamResp.Choices[0].Delta.Content
					fmt.Print(content)
					if streamResp.Choices[0].FinishReason != nil && *streamResp.Choices[0].FinishReason == "stop" {
						break
					}
				}
			}
		},
	}
	rootCmd.AddCommand(gptCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
