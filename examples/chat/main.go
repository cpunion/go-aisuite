package main

import (
	"context"
	"fmt"

	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/client"
)

func main() {
	// Initialize client with API keys
	c := client.New(&client.APIKey{
		// Set your OpenAI API key or leave empty to use environment variable OPENAI_API_KEY
		OpenAI: "",
		// Set your Anthropic API key or leave empty to use environment variable ANTHROPIC_API_KEY
		Anthropic: "",
	})

	// Make a chat completion request
	resp, err := c.ChatCompletion(context.Background(), aisuite.ChatCompletionRequest{
		Model: "openai:gpt-4o-mini", // or "anthropic:claude-3-5-haiku-20241022"
		Messages: []aisuite.ChatCompletionMessage{
			{
				Role:    aisuite.RoleUser,
				Content: "Hello, how are you?",
			},
		},
		MaxTokens: 10,
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}
