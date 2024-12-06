# go-aisuite

A cross-platform Go library for interacting with multiple AI providers' APIs, inspired by [aisuite](https://github.com/andrewyng/aisuite). Currently supports OpenAI and Anthropic providers with a unified interface.

## Features

- Unified interface for multiple AI providers
- Currently supports:
  - OpenAI (via [go-openai](https://github.com/sashabaranov/go-openai))
  - Anthropic (via [go-anthropic](https://github.com/liushuangls/go-anthropic))
- Carefully designed API that follows each provider's best practices
- Gradual and thoughtful addition of necessary interfaces and fields

## Installation

```bash
go get github.com/cpunion/go-aisuite
```

## Quick Start

See complete examples in the [examples](./examples) directory.

Basic usage:

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
	// Initialize client with API keys
	c := client.New(&client.APIKey{
		// Set your OpenAI API key or leave empty to use environment variable OPENAI_API_KEY
		OpenAI: "your-openai-api-key",
		// Set your Anthropic API key or leave empty to use environment variable ANTHROPIC_API_KEY
		Anthropic: "your-anthropic-api-key",
	})

	// Make a chat completion request
	resp, err := c.ChatCompletion(context.Background(), aisuite.ChatCompletionRequest{
		Model: "openai:gpt-4o-mini", // or "anthropic:claude-3-5-haiku-20241022"
		Messages: []aisuite.ChatCompletionMessage{
			{
				Role:    aisuite.User,
				Content: "Hello, how are you?",
			},
		},
	})
	if err != nil {
		panic(err)
	}

	fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
}

```

## Contributing

We welcome contributions! Please feel free to submit a Pull Request. We are carefully expanding the API surface area to maintain compatibility and usability across different providers.

## License

MIT License

## Acknowledgments

This project is inspired by [aisuite](https://github.com/andrewyng/aisuite) and builds upon the excellent work of:
- [go-openai](https://github.com/sashabaranov/go-openai)
- [go-anthropic](https://github.com/liushuangls/go-anthropic)
