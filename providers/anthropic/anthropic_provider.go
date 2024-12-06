package anthropic

import (
	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/providers"
)

const Name = "anthropic"

func init() {
	providers.RegisterProvider(Name, Provider{})
}

type Provider struct {
}

func (p Provider) NewClient(apiKey string) aisuite.Client {
	return NewClient(apiKey)
}
