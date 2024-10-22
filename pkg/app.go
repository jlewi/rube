package pkg

import (
	"github.com/go-logr/zapr"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/jlewi/monogo/files"
	"github.com/pkg/errors"
	"github.com/sashabaranov/go-openai"
	"go.uber.org/zap"
	"os"
)

type App struct {
	otelShutdownFn func()
	client         *openai.Client
}

func (a *App) Run(httpPort int, honeycombApiKeyFile string, openaiApiKeyFile string) error {
	if err := a.setupLogging(); err != nil {
		return errors.Wrapf(err, "Failed to setup logging")
	}

	if err := a.SetupHoneycomb(honeycombApiKeyFile); err != nil {
		return errors.Wrapf(err, "Failed to setup Honeycomb")
	}

	client, err := NewClient(openaiApiKeyFile)
	if err != nil {
		return err
	}
	a.client = client
	return a.Serve(httpPort)
}

// Serve sets up and runs the server
// This is blocking
func (a *App) Serve(httpPort int) error {
	s, err := NewServer(httpPort, a.client)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return s.Run()
}

func (a *App) setupLogging() error {
	// Configure encoder for JSON format
	c := zap.NewProductionConfig()
	// Use the keys used by cloud logging
	// https://cloud.google.com/logging/docs/structured-logging
	c.EncoderConfig.LevelKey = "severity"
	c.EncoderConfig.TimeKey = "time"
	c.EncoderConfig.MessageKey = "message"
	// We attach the function key to the logs because that is useful for identifying the function that generated the log.
	c.EncoderConfig.FunctionKey = "function"

	l, err := c.Build()
	if err != nil {
		return errors.Wrap(err, "failed to build logger")
	}

	zap.ReplaceGlobals(l)
	return nil
}

// SetupHoneycomb configures OTEL to export metrics to Honeycomb
func (a *App) SetupHoneycomb(apiKeyFile string) error {
	log := zapr.NewLogger(zap.L())
	log.Info("Configuring Honeycomb")

	// https://docs.honeycomb.io/send-data/go/opentelemetry-sdk/
	key, err := files.Read(apiKeyFile)
	if err != nil {
		return errors.Wrapf(err, "Could not read secret: %v", apiKeyFile)
	}

	serviceName := "rube"

	// The environment variable OTEL_SERVICE_NAME is the default for the honeycomb dataset.
	// https://docs.honeycomb.io/getting-data-in/opentelemetry/go-distro/
	// This will default to unknown. We don't want to use "unknown" as the default value so we override it.
	if os.Getenv("OTEL_SERVICE_NAME") != "" {
		serviceName = os.Getenv("OTEL_SERVICE_NAME")
		log.Info("environment variable OTEL_SERVICE_NAME is set", "service", serviceName)
	}
	log.Info("Setting OTEL_SERVICE_NAME service name", "service", serviceName)

	opts := []otelconfig.Option{
		honeycomb.WithApiKey(string(key)),
	}

	// See https://docs.honeycomb.io/send-data/go/opentelemetry-sdk/
	headers := map[string]string{
		"x-honeycomb-team": string(key),
	}
	endpoint := "https://api.honeycomb.io:443"
	opts = append(opts, otelconfig.WithServiceName(serviceName), otelconfig.WithHeaders(headers), otelconfig.WithExporterEndpoint(endpoint))

	// Configure Honeycomb
	otelShutdown, err := otelconfig.ConfigureOpenTelemetry(opts...)
	if err != nil {
		return errors.Wrapf(err, "error setting up open telemetry")
	}
	a.otelShutdownFn = otelShutdown
	return nil
}
