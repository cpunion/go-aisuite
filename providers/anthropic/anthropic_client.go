package anthropic

import (
	"context"
	"os"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"github.com/anthropics/anthropic-sdk-go/packages/ssestream"
	"github.com/cpunion/go-aisuite"
)

const (
	defaultMaxTokens = 2048
)

type Client struct {
	client *anthropic.Client
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("ANTHROPIC_API_KEY")
		if token == "" {
			panic("ANTHROPIC_API_KEY not found in environment variables")
		}
	}
	return &Client{client: anthropic.NewClient(option.WithAPIKey(token))}
}

func (c *Client) ChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (*aisuite.ChatCompletionResponse, error) {
	messages := make([]anthropic.MessageParam, len(req.Messages))
	for i, msg := range req.Messages {
		if msg.Content != "" {
			messages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
		}
	}

	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	resp, err := c.client.Messages.New(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(req.Model)),
		Messages:  anthropic.F(messages),
		MaxTokens: anthropic.F(maxTokens),
	})
	if err != nil {
		return nil, err
	}

	content := ""
	if len(resp.Content) > 0 {
		content = resp.Content[0].Text
	}

	return &aisuite.ChatCompletionResponse{
		Choices: []aisuite.ChatCompletionChoice{
			{
				Message: aisuite.ChatCompletionMessage{
					Role:    aisuite.Role(string(resp.Role)),
					Content: content,
				},
			},
		},
	}, nil
}

func (c *Client) StreamChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (aisuite.ChatCompletionStream, error) {
	messages := make([]anthropic.MessageParam, len(req.Messages))
	for i, msg := range req.Messages {
		if msg.Content != "" {
			messages[i] = anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content))
		}
	}

	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	stream := c.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(req.Model)),
		Messages:  anthropic.F(messages),
		MaxTokens: anthropic.F(maxTokens),
	})

	return &chatCompletionStream{
		stream: stream,
	}, nil
}

type chatCompletionStream struct {
	stream *ssestream.Stream[anthropic.MessageStreamEvent]
}

func (s *chatCompletionStream) Recv() (aisuite.ChatCompletionStreamResponse, error) {
	for {
		if !s.stream.Next() {
			continue
		}
		if err := s.stream.Err(); err != nil {
			return aisuite.ChatCompletionStreamResponse{}, err
		}

		event := s.stream.Current()

		switch event.Type {
		case anthropic.MessageStreamEventTypeMessageDelta:
			event := event.Delta.(anthropic.MessageDeltaEventDelta)
			if event.StopReason != "" {
				return aisuite.ChatCompletionStreamResponse{
					Choices: []aisuite.ChatCompletionStreamChoice{
						{
							FinishReason: string(event.StopReason),
						},
					},
				}, nil
			}
		case anthropic.MessageStreamEventTypeContentBlockDelta:
			delta := event.Delta.(anthropic.ContentBlockDeltaEventDelta)
			return aisuite.ChatCompletionStreamResponse{
				Choices: []aisuite.ChatCompletionStreamChoice{
					{
						Delta: aisuite.ChatCompletionStreamChoiceDelta{
							Role:    "assistant",
							Content: delta.Text,
						},
					},
				},
			}, nil
		}
	}
}

func (s *chatCompletionStream) Close() error {
	return s.stream.Close()
}
