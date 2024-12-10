package client

import (
	"context"
	"fmt"
	"strings"

	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/providers"
	"github.com/cpunion/go-aisuite/providers/anthropic"
	"github.com/cpunion/go-aisuite/providers/gemini"
	"github.com/cpunion/go-aisuite/providers/groq"
	"github.com/cpunion/go-aisuite/providers/openai"
	"github.com/cpunion/go-aisuite/providers/sambanova"
)

const (
	ErrUnknownProvider = "unknown provider"
)

type APIKey struct {
	OpenAI    string
	Anthropic string
	Sambanova string
	Gemini    string
	Groq      string
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
	opts := providers.Options{}
	switch providerName {
	case openai.Name:
		opts.Token = c.apiKey.OpenAI
	case anthropic.Name:
		opts.Token = c.apiKey.Anthropic
	case gemini.Name:
		opts.BaseURL = "https://generativelanguage.googleapis.com/v1beta/openai/"
		opts.Token = c.apiKey.Gemini
	case sambanova.Name:
		opts.BaseURL = "https://api.sambanova.ai/v1/"
		opts.Token = c.apiKey.Sambanova
	case groq.Name:
		opts.Token = c.apiKey.Groq
		opts.BaseURL = "https://api.groq.com/openai/v1/"
	}
	return provider.NewClient(opts), toks[1]
}
