package aisuite

import "context"

type Client interface {
	ChatCompletion(ctx context.Context, request ChatCompletionRequest) (*ChatCompletionResponse, error)
	StreamChatCompletion(ctx context.Context, request ChatCompletionRequest) (ChatCompletionStream, error)
}
