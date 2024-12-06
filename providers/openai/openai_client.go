package openai

import (
	"context"
	"log/slog"
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
			Role:    toOpenAIRole(msg.Role),
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
				Role:    fromOpenAIRole(choice.Message.Role),
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
	var role aisuite.Role
	for i, choice := range resp.Choices {
		if choice.Delta.Role != "" {
			// Just use first role, other roles are blank
			role = fromOpenAIRole(choice.Delta.Role)
		}
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
				Role:         role,
				FunctionCall: funcCall,
				ToolCalls:    toolCalls,
				Refusal:      choice.Delta.Refusal,
			},
			FinishReason: fromOpenAIFinishReason(choice.FinishReason),
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
			Role:    toOpenAIRole(msg.Role),
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

func fromOpenAIFinishReason(reason ai.FinishReason) aisuite.FinishReason {
	switch reason {
	case "":
		return aisuite.FinishReasonNone
	case ai.FinishReasonStop:
		return aisuite.FinishReasonStop
	case ai.FinishReasonLength:
		return aisuite.FinishReasonMaxTokens
	}
	return aisuite.FinishReason("unknown: " + string(reason))
}

func fromOpenAIRole(role string) aisuite.Role {
	switch role {
	case "user":
		return aisuite.RoleUser
	case "system":
		return aisuite.RoleSystem
	case "assistant":
		return aisuite.RoleAssistant
	}
	slog.Warn("unknown openai role, should handle this", "role", role)
	return aisuite.Role(string(role))
}

func toOpenAIRole(role aisuite.Role) string {
	switch role {
	case aisuite.RoleUser:
		return "user"
	case aisuite.RoleSystem:
		return "system"
	case aisuite.RoleAssistant:
		return "assistant"
	}
	slog.Warn("can't convert aisuite role to openai role, should handle this", "role", role)
	return string(role)
}
