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
		OpenAI:    "", // Set your OpenAI API key or use OPENAI_API_KEY env
		Anthropic: "", // Set your Anthropic API key or use ANTHROPIC_API_KEY env
	})

	// Create a streaming chat completion request
	stream, err := c.StreamChatCompletion(context.Background(), aisuite.ChatCompletionRequest{
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
	defer stream.Close()

	// Read the response stream
	for {
		resp, err := stream.Recv()
		if err != nil {
			panic(err)
		}
		if len(resp.Choices) == 0 {
			fmt.Println("No choices")
			break
		}
		if resp.Choices[0].FinishReason != "" {
			fmt.Printf("\nStream finished: %s\n", resp.Choices[0].FinishReason)
			break
		}
		fmt.Print(resp.Choices[0].Delta.Content)
	}
}
