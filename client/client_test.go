package client

import (
	"context"
	"testing"
	"time"

	"github.com/cpunion/go-aisuite"
)

func TestChatCompletion(t *testing.T) {
	client := New(nil)
	models := []string{
		"openai:gpt-4o-mini",
		"anthropic:claude-3-5-haiku-20241022",
	}
	for _, model := range models {
		resp, err := client.ChatCompletion(context.Background(), aisuite.ChatCompletionRequest{
			Model: model,
			Messages: []aisuite.ChatCompletionMessage{
				{
					Role:    aisuite.User,
					Content: "Hello",
				},
			},
			MaxTokens: 10,
		})
		if err != nil {
			t.Fatal(err)
		}
		t.Logf("resp: %#v", resp)
	}
}

func TestStreamChatCompletion(t *testing.T) {
	client := New(nil)
	models := []string{
		"openai:gpt-4o-mini",
		"anthropic:claude-3-5-haiku-20241022",
	}
	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()

			stream, err := client.StreamChatCompletion(ctx, aisuite.ChatCompletionRequest{
				Model: model,
				Messages: []aisuite.ChatCompletionMessage{
					{
						Role:    aisuite.User,
						Content: "Hello",
					},
				},
				MaxTokens: 10,
			})
			if err != nil {
				t.Fatal(err)
			}
			defer stream.Close()

			// Read all chunks from the stream
			for {
				select {
				case <-ctx.Done():
					t.Fatal("context deadline exceeded")
					return
				default:
					chunk, err := stream.Recv()
					if err != nil {
						if err.Error() == "EOF" {
							return
						}
						t.Fatal(err)
					}

					t.Logf("chunk: %#v", chunk)
					if len(chunk.Choices) > 0 && chunk.Choices[0].FinishReason != "" {
						t.Logf("finish reason: %s", chunk.Choices[0].FinishReason)
						return
					}
				}
			}
		})
	}
}
