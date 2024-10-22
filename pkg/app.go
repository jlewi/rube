package pkg

import (
	"github.com/go-logr/zapr"
	"github.com/honeycombio/honeycomb-opentelemetry-go"
	"github.com/honeycombio/otel-config-go/otelconfig"
	"github.com/jlewi/monogo/files"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"os"
)

type App struct {
	otelShutdownFn func()
}

func (a *App) Run(httpPort int, honeycombApiKeyFile string) error {
	if err := a.SetupHoneycomb(honeycombApiKeyFile); err != nil {
		return errors.Wrapf(err, "Failed to setup Honeycomb")
	}

	return a.Serve(httpPort)
}

// Serve sets up and runs the server
// This is blocking
func (a *App) Serve(httpPort int) error {

	s, err := NewServer(httpPort)

	if err != nil {
		return err
	}

	if err != nil {
		return err
	}

	return s.Run()
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
