package main

import (
	"fmt"
	"github.com/go-logr/zapr"
	"github.com/jlewi/foyle/rube/pkg"
	"go.uber.org/zap"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func main() {
	var honeycombAPIKeyFile string
	var openaiAPIKeyFile string
	var httpPort int
	var rootCmd = &cobra.Command{
		Use:   "rube",
		Short: "A Gin HTTP server instrumented with OTEL, logging with logr/zapr",
		Run: func(cmd *cobra.Command, args []string) {
			app := &pkg.App{}
			err := app.Run(httpPort, honeycombAPIKeyFile, openaiAPIKeyFile)

			if err != nil {
				log := zapr.NewLogger(zap.L())
				log.Error(err, "Failed to run the application")
				fmt.Fprintf(os.Stdout, "Error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	rootCmd.PersistentFlags().StringVarP(&honeycombAPIKeyFile, "honeycomb-apikey", "k", "", "Path to the Honeycomb API key file")
	rootCmd.PersistentFlags().StringVarP(&openaiAPIKeyFile, "openai-apikey", "", "", "Path to the OpenAI API key file")
	rootCmd.PersistentFlags().IntVarP(&httpPort, "port", "p", 8080, "Port to serve on")
	rootCmd.MarkFlagRequired("honeycomb-apikey")
	rootCmd.MarkFlagRequired("openai-apikey")

	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to execute root command: %v", err)
	}
}
