package providers

import "github.com/cpunion/go-aisuite"

type Provider interface {
	NewClient(apiKey string) aisuite.Client
}

var providers = make(map[string]Provider)

func RegisterProvider(name string, provider Provider) {
	providers[name] = provider
}

func GetProvider(name string) (Provider, bool) {
	provider, ok := providers[name]
	return provider, ok
}
