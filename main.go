package main

import (
	"fmt"
	"github.com/jlewi/foyle/rube/pkg"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var honeycombAPIKeyFile string
	var httpPort int
	var rootCmd = &cobra.Command{
		Use:   "otel-gin-server",
		Short: "A Gin HTTP server instrumented with OTEL, logging with logr/zapr",
		Run: func(cmd *cobra.Command, args []string) {
			app := &pkg.App{}
			err := app.Run(httpPort, honeycombAPIKeyFile)

			if err != nil {
				fmt.Fprintf(os.Stdout, "Error: %v\n", err)
				os.Exit(1)
			}

			//tp, err := setupTracer(honeycombAPIKey)
			//if err != nil {
			//	logger.Error(err, "Failed to initialize tracer")
			//	os.Exit(1)
			//}
			//defer func() {
			//	if err := tp.Shutdown(context.Background()); err != nil {
			//		logger.Error(err, "Failed to shutdown tracer")
			//	}
			//}()
		},
	}

	rootCmd.PersistentFlags().StringVarP(&honeycombAPIKeyFile, "honeycomb-apikey", "k", "", "Path to the Honeycomb API key file")
	rootCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 8080, "Port to serve on")
	rootCmd.MarkPersistentFlagRequired("apikey-file")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}
