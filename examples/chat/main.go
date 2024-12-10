package main

import (
	"context"
	"fmt"

	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/client"
)

func main() {
	// Initialize client with environment variables
	c := client.New(nil)

	// Or initialize client with API keys
	// c := client.New(&client.APIKey{
	// 	OpenAI:    "", // Set your OpenAI API key or keep empty to use OPENAI_API_KEY env
	// 	Anthropic: "", // Set your Anthropic API key or keep empty to use ANTHROPIC_API_KEY env
	// 	Groq:      "", // Set your Groq API key or keep empty to use GROQ_API_KEY env
	// 	Gemini:    "", // Set your Gemini API key or keep empty to use GEMINI_API_KEY env
	// 	Sambanova: "", // Set your SambaNova API key or keep empty to use SAMBANOVA_API_KEY env
	// })

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
