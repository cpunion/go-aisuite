package client

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/cpunion/go-aisuite"
)

var testModels = []string{
	"openai:gpt-4o-mini",
	"anthropic:claude-3-5-haiku-20241022",
	"gemini:gemini-1.5-flash-latest",
	"sambanova:Meta-Llama-3.2-1B-Instruct",
	"groq:llama-3.1-8b-instant",
}

type testCase struct {
	name             string
	model            string
	prompt           string
	maxTokens        int
	wantFinishReason aisuite.FinishReason
}

func generateTestCases() []testCase {
	var cases []testCase

	// Test cases for max_tokens finish reason
	longStory := "Tell me a very long story about a magical adventure with dragons, wizards, and epic battles."
	for _, model := range testModels {
		cases = append(cases, testCase{
			name:             fmt.Sprintf("%s_maxtoken_test", strings.Split(model, ":")[1]),
			model:            model,
			prompt:           longStory,
			maxTokens:        5,
			wantFinishReason: aisuite.FinishReasonMaxTokens,
		})
	}

	// Test cases for normal stop
	for _, model := range testModels {
		cases = append(cases, testCase{
			name:             fmt.Sprintf("%s_normal_stop", strings.Split(model, ":")[1]),
			model:            model,
			prompt:           "Hi (shortly respond)",
			maxTokens:        20,
			wantFinishReason: aisuite.FinishReasonStop,
		})
	}

	return cases
}

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
	models := testModels
	for _, model := range models {
		t.Run(model, func(t *testing.T) {
			wd, _ := os.Getwd()
			t.Logf("Working directory: %s", wd)
			withTimeout(t, 10*time.Second, func(ctx context.Context) {
				resp, err := client.ChatCompletion(ctx, aisuite.ChatCompletionRequest{
					Model: model,
					Messages: []aisuite.ChatCompletionMessage{
						{
							Role:    aisuite.RoleUser,
							Content: "Hi",
						},
					},
					MaxTokens: 30,
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
	cases := generateTestCases()

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withTimeout(t, 10*time.Second, func(ctx context.Context) {
				stream, err := client.StreamChatCompletion(ctx, aisuite.ChatCompletionRequest{
					Model: tc.model,
					Messages: []aisuite.ChatCompletionMessage{
						{
							Role:    aisuite.RoleUser,
							Content: tc.prompt,
						},
					},
					MaxTokens: tc.maxTokens,
				})
				if err != nil {
					t.Fatal(err)
				}
				defer stream.Close()

				var content string
				var finishReason aisuite.FinishReason
				for {
					resp, err := stream.Recv()
					if err != nil {
						t.Fatal(err)
					}
					if len(resp.Choices) == 0 {
						fmt.Println("No choices")
						break
					}
					if resp.Choices[0].FinishReason != "" {
						finishReason = resp.Choices[0].FinishReason
						break
					}
					content += resp.Choices[0].Delta.Content
				}

				if finishReason != tc.wantFinishReason {
					t.Errorf("got finish reason %q, want %q", finishReason, tc.wantFinishReason)
				}

				fmt.Printf("Test case %s:\nPrompt: %s\nResponse: %s\nFinish reason: %s\n\n",
					tc.name, tc.prompt, content, finishReason)
			})
		})
	}
}
