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
		OpenAI:    "", // Set your OpenAI API key or keep empty to use OPENAI_API_KEY env
		Anthropic: "", // Set your Anthropic API key or keep empty to use ANTHROPIC_API_KEY env
		Groq:      "", // Set your Groq API key or keep empty to use GROQ_API_KEY env
		Gemini:    "", // Set your Gemini API key or keep empty to use GEMINI_API_KEY env
		Sambanova: "", // Set your SambaNova API key or keep empty to use SAMBANOVA_API_KEY env
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
