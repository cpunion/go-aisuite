package anthropic

import (
	"context"
	"log/slog"
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
					Role:    fromAnthropicRole(resp.Role),
					Content: content,
				},
			},
		},
	}, nil
}

func (c *Client) StreamChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (aisuite.ChatCompletionStream, error) {
	system := make([]anthropic.TextBlockParam, 0, 1)
	messages := make([]anthropic.MessageParam, 0, len(req.Messages))
	for _, msg := range req.Messages {
		switch msg.Role {
		case aisuite.RoleSystem:
			system = append(system, anthropic.NewTextBlock(msg.Content))
		case aisuite.RoleUser:
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
		case aisuite.RoleAssistant:
			messages = append(messages, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.Content)))
		default:
			messages = append(messages, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.Content)))
		}
	}

	maxTokens := int64(req.MaxTokens)
	if maxTokens == 0 {
		maxTokens = defaultMaxTokens
	}

	stream := c.client.Messages.NewStreaming(ctx, anthropic.MessageNewParams{
		Model:     anthropic.F(anthropic.Model(req.Model)),
		System:    anthropic.F(system),
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
							FinishReason: fromAnthropicStopReason(event.StopReason),
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
							Role:    aisuite.RoleAssistant,
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

func fromAnthropicStopReason(stopReason anthropic.MessageDeltaEventDeltaStopReason) aisuite.FinishReason {
	switch stopReason {
	case "":
		return aisuite.FinishReasonNone
	case "end_turn":
		return aisuite.FinishReasonStop
	case "max_tokens":
		return aisuite.FinishReasonMaxTokens
	case "content_filter":
		return aisuite.FinishReasonContentFilter
	default:
		return aisuite.FinishReason(string(aisuite.FinishReasonUnknown) + "(" + string(stopReason) + ")")
	}
}

func fromAnthropicRole(role anthropic.MessageRole) aisuite.Role {
	switch role {
	case anthropic.MessageRoleAssistant:
		return aisuite.RoleAssistant
	default:
		slog.Warn("can't convert anthropic role to aisuite role, should handle this", "role", role)
		return aisuite.Role(string(role))
	}
}
