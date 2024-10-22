package pkg

import (
	"github.com/go-logr/zapr"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/jlewi/monogo/files"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.uber.org/zap"
)

// NewClient helper function to create a new OpenAI client
func NewClient(apiKeyFile string) (*openai.Client, error) {
	log := zapr.NewLogger(zap.L())
	// ************************************************************************
	// Setup middleware
	// ************************************************************************

	// Handle retryable errors
	// To handle retryable errors we use hashi corp's retryable client. This client will automatically retry on
	// retryable errors like 429; rate limiting
	retryClient := retryablehttp.NewClient()
	httpClient := retryClient.StandardClient()
	httpClient.Transport = otelhttp.NewTransport(httpClient.Transport)

	log.Info("Configuring OpenAI client")

	apiKey, err := files.Read(apiKeyFile)
	if err != nil {
		return nil, errors.Wrapf(err, "Could not read API key from file: %s", apiKeyFile)
	}

	clientConfig := openai.DefaultConfig(string(apiKey))

	clientConfig.HTTPClient = httpClient
	client := openai.NewClientWithConfig(clientConfig)

	return client, nil
}
