package openai

import (
	"context"
	"os"

	"github.com/cpunion/go-aisuite"
	ai "github.com/sashabaranov/go-openai"
)

type Client struct {
	client *ai.Client
}

func NewClient(token string) *Client {
	if token == "" {
		token = os.Getenv("OPENAI_API_KEY")
		if token == "" {
			panic("OPENAI_API_KEY not found in environment variables")
		}
	}
	return &Client{client: ai.NewClient(token)}
}

func (c *Client) ChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (*aisuite.ChatCompletionResponse, error) {
	aiMessages := make([]ai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		aiMessages[i] = ai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	chatReq := ai.ChatCompletionRequest{
		Model:     req.Model,
		MaxTokens: req.MaxTokens,
		Stream:    req.Stream,
		Messages:  aiMessages,
	}
	resp, err := c.client.CreateChatCompletion(ctx, chatReq)
	if err != nil {
		return nil, err
	}
	choices := make([]aisuite.ChatCompletionChoice, len(resp.Choices))
	for i, choice := range resp.Choices {
		choices[i] = aisuite.ChatCompletionChoice{
			Message: aisuite.ChatCompletionMessage{
				Content: choice.Message.Content,
			},
		}
	}
	return &aisuite.ChatCompletionResponse{Choices: choices}, nil
}

type chatCompletionStream struct {
	stream *ai.ChatCompletionStream
}

func (c *chatCompletionStream) Recv() (aisuite.ChatCompletionStreamResponse, error) {
	resp, err := c.stream.Recv()
	if err != nil {
		return aisuite.ChatCompletionStreamResponse{}, err
	}
	choices := make([]aisuite.ChatCompletionStreamChoice, len(resp.Choices))
	for i, choice := range resp.Choices {
		var funcCall *aisuite.FunctionCall
		if choice.Delta.FunctionCall != nil {
			funcCall = &aisuite.FunctionCall{
				Name: choice.Delta.FunctionCall.Name,
				Args: choice.Delta.FunctionCall.Arguments,
			}
		}
		toolCalls := make([]aisuite.ToolCall, len(choice.Delta.ToolCalls))
		for j, toolCall := range choice.Delta.ToolCalls {
			toolCalls[j] = aisuite.ToolCall{
				ID:   toolCall.ID,
				Tool: string(toolCall.Type),
				Function: aisuite.FunctionCall{
					Name: toolCall.Function.Name,
					Args: toolCall.Function.Arguments,
				},
			}
		}
		choices[i] = aisuite.ChatCompletionStreamChoice{
			Delta: aisuite.ChatCompletionStreamChoiceDelta{
				Content:      choice.Delta.Content,
				Role:         choice.Delta.Role,
				FunctionCall: funcCall,
				ToolCalls:    toolCalls,
				Refusal:      choice.Delta.Refusal,
			},
			FinishReason: string(choice.FinishReason),
		}
	}
	return aisuite.ChatCompletionStreamResponse{Choices: choices}, nil
}

func (c *chatCompletionStream) Close() error {
	return c.stream.Close()
}

func (c *Client) StreamChatCompletion(ctx context.Context, req aisuite.ChatCompletionRequest) (aisuite.ChatCompletionStream, error) {
	aiMessages := make([]ai.ChatCompletionMessage, len(req.Messages))
	for i, msg := range req.Messages {
		aiMessages[i] = ai.ChatCompletionMessage{
			Role:    string(msg.Role),
			Content: msg.Content,
		}
	}
	chatReq := ai.ChatCompletionRequest{
		Model:     req.Model,
		Messages:  aiMessages,
		MaxTokens: req.MaxTokens,
		Stream:    true,
	}
	s, err := c.client.CreateChatCompletionStream(ctx, chatReq)
	if err != nil {
		return nil, err
	}
	return &chatCompletionStream{stream: s}, nil
}
