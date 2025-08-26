package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/eggybyte-technology/yao-portal/internal/config"
	"github.com/eggybyte-technology/yao-portal/internal/server"
)

const (
	// DefaultConfigPath is the default path for configuration file
	DefaultConfigPath = "configs/portal.yaml"

	// AppName is the application name
	AppName = "YaoPortal"

	// AppVersion is the application version
	AppVersion = "v1.0.0"
)

var (
	configPath = flag.String("config", DefaultConfigPath, "Path to configuration file")
	version    = flag.Bool("version", false, "Print version information")
	help       = flag.Bool("help", false, "Print help information")
)

func main() {
	flag.Parse()

	// Print version information
	if *version {
		log.Printf("%s %s", AppName, AppVersion)
		os.Exit(0)
	}

	// Print help information
	if *help {
		printHelp()
		os.Exit(0)
	}

	// Load configuration
	cfg, err := config.LoadConfig(*configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Create server instance
	srv, err := server.NewPortalServer(cfg)
	if err != nil {
		log.Fatalf("Failed to create portal server: %v", err)
	}

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Starting %s %s...", AppName, AppVersion)
		if err := srv.Start(ctx); err != nil {
			log.Fatalf("Failed to start portal server: %v", err)
		}
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping server gracefully...")

	// Cancel context to trigger shutdown
	cancel()

	// Stop server
	if err := srv.Stop(ctx); err != nil {
		log.Fatalf("Failed to stop portal server gracefully: %v", err)
	}

	log.Println("Server stopped successfully")
}

// printHelp prints help information
func printHelp() {
	log.Printf(`%s %s - Ethereum-compatible JSON-RPC API Gateway for YaoVerse

Usage: %s [options]

Options:
  -config string
        Path to configuration file (default "%s")
  -version
        Print version information
  -help
        Print this help information

Environment Variables:
  YAO_PORTAL_CONFIG    Configuration file path (overrides -config flag)
  YAO_PORTAL_LOG_LEVEL Log level (debug, info, warn, error)

Examples:
  %s                           # Use default configuration
  %s -config /etc/portal.yaml  # Use custom configuration file
  %s -version                  # Print version information

For more information, visit: https://github.com/eggybyte-technology/yao-verse
`, AppName, AppVersion, os.Args[0], DefaultConfigPath, os.Args[0], os.Args[0], os.Args[0])
}
