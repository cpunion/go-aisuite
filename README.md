# go-aisuite

A cross-platform Go library for interacting with multiple AI providers' APIs, inspired by [aisuite](https://github.com/andrewyng/aisuite). Currently supports OpenAI and Anthropic providers with a unified interface.


[![Build Status](https://github.com/cpunion/go-aisuite/actions/workflows/test.yml/badge.svg)](https://github.com/cpunion/go-aisuite/actions/workflows/test.yml)
[![codecov](https://codecov.io/github/cpunion/go-aisuite/graph/badge.svg?token=uATQa0RzPL)](https://codecov.io/github/cpunion/go-aisuite)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/cpunion/go-aisuite)
[![GitHub commits](https://badgen.net/github/commits/cpunion/go-aisuite)](https://GitHub.com/Naereen/cpunion/go-aisuite/commit/)
[![GitHub release](https://img.shields.io/github/v/tag/cpunion/go-aisuite.svg?label=release)](https://github.com/cpunion/go-aisuite/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/cpunion/go-aisuite)](https://goreportcard.com/report/github.com/cpunion/go-aisuite)
[![Go Reference](https://pkg.go.dev/badge/github.com/cpunion/go-aisuite.svg)](https://pkg.go.dev/github.com/cpunion/go-aisuite)

## Features

- Unified interface for multiple AI providers
- Currently supports:
  - OpenAI (via [go-openai](https://github.com/sashabaranov/go-openai))
  - Anthropic (via [official SDK](https://github.com/anthropics/anthropic-sdk-go))
  - Groq (via OpenAI-compatible API)
  - Gemini (via OpenAI-compatible API)
  - SambaNova (via OpenAI-compatible API)
- Carefully designed API that follows each provider's best practices
- Gradual and thoughtful addition of necessary interfaces and fields

## Installation

```bash
go get github.com/cpunion/go-aisuite
```

## Quick Start

See complete examples in the [examples](./examples) directory.

### Chat

<!-- embedme examples/chat/main.go -->

```go
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

```

### Stream

<!-- embedme examples/stream/main.go -->

```go
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

```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. We are carefully expanding the API surface area to maintain compatibility and usability across different providers.

## License

MIT License

## Acknowledgments

This project is inspired by [aisuite](https://github.com/andrewyng/aisuite) and builds upon the excellent work of:
- [go-openai](https://github.com/sashabaranov/go-openai)
- [anthropic-sdk-go](https://github.com/anthropics/anthropic-sdk-go)
