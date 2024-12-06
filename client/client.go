package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/providers"
	_ "github.com/cpunion/go-aisuite/providers/anthropic"
	_ "github.com/cpunion/go-aisuite/providers/openai"
)

const (
	ErrUnknownProvider = "unknown provider"
)

type APIKey struct {
	OpenAI    string
	Anthropic string
}

type AdaptiveClient struct {
	apiKey *APIKey
}

func New(apiKey *APIKey) aisuite.Client {
	if apiKey == nil {
		apiKey = &APIKey{}
	}
	return AdaptiveClient{apiKey: apiKey}
}

func (c AdaptiveClient) ChatCompletion(ctx context.Context, request aisuite.ChatCompletionRequest) (*aisuite.ChatCompletionResponse, error) {
	client, model := c.getClientAndModel(request.Model)
	newReq := request
	newReq.Model = model
	return client.ChatCompletion(ctx, newReq)
}

func (c AdaptiveClient) StreamChatCompletion(ctx context.Context, request aisuite.ChatCompletionRequest) (aisuite.ChatCompletionStream, error) {
	client, model := c.getClientAndModel(request.Model)
	newReq := request
	newReq.Model = model
	return client.StreamChatCompletion(ctx, newReq)
}

func (c AdaptiveClient) getClientAndModel(model string) (aisuite.Client, string) {
	toks := strings.SplitN(model, ":", 2)
	providerName := toks[0]
	provider, ok := providers.GetProvider(providerName)
	if !ok {
		panic(fmt.Sprintf("%s: %s", ErrUnknownProvider, providerName))
	}
	var apiKey string
	switch providerName {
	case "openai":
		apiKey = c.apiKey.OpenAI
	case "anthropic":
		apiKey = c.apiKey.Anthropic
	}
	return provider.NewClient(apiKey), toks[1]
}
