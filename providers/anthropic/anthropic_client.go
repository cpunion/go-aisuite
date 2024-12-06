package anthropic

import (
	"context"
	"os"

	"github.com/cpunion/go-aisuite"
	anthropicapi "github.com/liushuangls/go-anthropic/v2"
)

const (
	defaultMaxTokens = 2048
)

type Client struct {
	client *anthropicapi.Client
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("ANTHROPIC_API_KEY")
		if token == "" {
			panic("ANTHROPIC_API_KEY not found in environment variables")
		}
	}
	return &Client{client: anthropicapi.NewClient(token)}
}

func (c *Client) ChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (*aisuite.ChatCompletionResponse, error) {
	messages := make([]anthropicapi.Message, len(req.Messages))
	for i, msg := range req.Messages {
		var content []anthropicapi.MessageContent
		if msg.Content != "" {
			content = []anthropicapi.MessageContent{
				anthropicapi.NewTextMessageContent(msg.Content),
			}
		}
		messages[i] = anthropicapi.Message{
			Role:    anthropicapi.ChatRole(msg.Role),
			Content: content,
		}
	}

	request := anthropicapi.MessagesRequest{
		Model:     anthropicapi.Model(req.Model),
		Messages:  messages,
		MaxTokens: req.MaxTokens,
	}
	if request.MaxTokens == 0 {
		request.MaxTokens = defaultMaxTokens
	}

	resp, err := c.client.CreateMessages(ctx, request)
	if err != nil {
		return nil, err
	}

	content := ""
	if len(resp.Content) > 0 {
		content = resp.Content[0].GetText()
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

type chatCompletionStream struct {
	ctx    context.Context
	client *anthropicapi.Client
	req    anthropicapi.MessagesStreamRequest
}

func (s *chatCompletionStream) Recv() (aisuite.ChatCompletionStreamResponse, error) {
	resp, err := s.client.CreateMessagesStream(s.ctx, s.req)
	if err != nil {
		return aisuite.ChatCompletionStreamResponse{}, err
	}

	// Extract content from the response
	content := ""
	if len(resp.Content) > 0 {
		content = resp.Content[0].GetText()
	}

	return aisuite.ChatCompletionStreamResponse{
		Choices: []aisuite.ChatCompletionStreamChoice{
			{
				Delta: aisuite.ChatCompletionStreamChoiceDelta{
					Role:    string(resp.Role),
					Content: content,
				},
				FinishReason: string(resp.StopReason),
			},
		},
	}, nil
}

func (s *chatCompletionStream) Close() error {
	return nil
}

func (c *Client) StreamChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (aisuite.ChatCompletionStream, error) {
	messages := make([]anthropicapi.Message, len(req.Messages))
	for i, msg := range req.Messages {
		var content []anthropicapi.MessageContent
		if msg.Content != "" {
			content = []anthropicapi.MessageContent{
				anthropicapi.NewTextMessageContent(msg.Content),
			}
		}
		messages[i] = anthropicapi.Message{
			Role:    anthropicapi.ChatRole(msg.Role),
			Content: content,
		}
	}

	streamReq := anthropicapi.MessagesStreamRequest{
		MessagesRequest: anthropicapi.MessagesRequest{
			Model:     anthropicapi.Model(req.Model),
			Messages:  messages,
			MaxTokens: req.MaxTokens,
			Stream:    true,
		},
	}

	return &chatCompletionStream{
		ctx:    ctx,
		client: c.client,
		req:    streamReq,
	}, nil
}
