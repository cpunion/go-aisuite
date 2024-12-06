package client

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/cpunion/go-aisuite"
)

func withTimeout(t *testing.T, timeout time.Duration, fn func(ctx context.Context)) {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	done := make(chan struct{})
	go func() {
		fn(ctx)
		close(done)
	}()

	select {
	case <-done:
		return
	case <-ctx.Done():
		t.Fatal("test timeout")
	}
}

func TestChatCompletion(t *testing.T) {
	client := New(nil)
	models := []string{
		"openai:gpt-4o-mini",
		"anthropic:claude-3-5-haiku-20241022",
	}
	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			withTimeout(t, 30*time.Second, func(ctx context.Context) {
				resp, err := client.ChatCompletion(ctx, aisuite.ChatCompletionRequest{
					Model: model,
					Messages: []aisuite.ChatCompletionMessage{
						{
							Role:    "user",
							Content: "Hello",
						},
					},
				})
				if err != nil {
					t.Fatal(err)
				}
				if len(resp.Choices) == 0 {
					t.Fatal("no choices")
				}
				fmt.Printf("Response: %s\n", resp.Choices[0].Message.Content)
			})
		})
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
			withTimeout(t, 30*time.Second, func(ctx context.Context) {
				stream, err := client.StreamChatCompletion(ctx, aisuite.ChatCompletionRequest{
					Model: model,
					Messages: []aisuite.ChatCompletionMessage{
						{
							Role:    "user",
							Content: "Hello",
						},
					},
					MaxTokens: 10,
				})
				if err != nil {
					t.Fatal(err)
				}
				defer stream.Close()

				var content string
				for {
					resp, err := stream.Recv()
					if err != nil {
						t.Fatal(err)
					}
					if len(resp.Choices) == 0 {
						break
					}
					if resp.Choices[0].FinishReason != "" {
						fmt.Printf("Stream stop reason: %s\n", resp.Choices[0].FinishReason)
						break
					}
					fmt.Printf("Stream Response: %s\n", resp.Choices[0].Delta.Content)
					content += resp.Choices[0].Delta.Content
				}
				fmt.Printf("Stream Response: %s\n", content)
			})
		})
	}
}
