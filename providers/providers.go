package providers

import "github.com/cpunion/go-aisuite"

type Options struct {
	BaseURL string
	Token   string
}

type Option func(o Options) Options

func WithToken(token string) Option {
	return func(o Options) Options {
		o.Token = token
		return o
	}
}

func WithBaseURL(baseURL string) Option {
	return func(o Options) Options {
		o.BaseURL = baseURL
		return o
	}
}

type Provider interface {
	NewClient(options Options) aisuite.Client
}

var providers = make(map[string]Provider)

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

func GetProvider(name string) (Provider, bool) {
	provider, ok := providers[name]
	return provider, ok
}
