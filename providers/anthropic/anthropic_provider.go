package anthropic

import (
	"os"

	"github.com/cpunion/go-aisuite"
	"github.com/cpunion/go-aisuite/providers"
)

const Name = "anthropic"
const apiKeyEnvVar = "ANTHROPIC_API_KEY"

func init() {
	providers.RegisterProvider(Name, Provider{})
}

type Provider struct {
}

func (p Provider) NewClient(opts providers.Options) aisuite.Client {
	if opts.Token == "" {
		opts.Token = os.Getenv(apiKeyEnvVar)
		if opts.Token == "" {
			panic(apiKeyEnvVar + " not found in environment variables")
		}
	}
	return NewClient(opts)
}
