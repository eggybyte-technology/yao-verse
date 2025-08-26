package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	// DefaultConfigPath is the default path for configuration file
	DefaultConfigPath = "configs/archive.yaml"

	// AppName is the application name
	AppName = "YaoArchive"

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

	// TODO: Load configuration
	// cfg, err := config.LoadConfig(*configPath)
	// if err != nil {
	//     log.Fatalf("Failed to load configuration: %v", err)
	// }

	// TODO: Create server instance
	// srv, err := server.NewArchiveServer(cfg)
	// if err != nil {
	//     log.Fatalf("Failed to create archive server: %v", err)
	// }

	// Create context for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Set up signal handling for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Start server
	go func() {
		log.Printf("Starting %s %s...", AppName, AppVersion)
		// TODO: Start actual server
		// if err := srv.Start(ctx); err != nil {
		//     log.Fatalf("Failed to start archive server: %v", err)
		// }

		// For now, just block until context is cancelled
		<-ctx.Done()
	}()

	// Wait for shutdown signal
	<-sigChan
	log.Println("Received shutdown signal, stopping server gracefully...")

	// Cancel context to trigger shutdown
	cancel()

	// TODO: Stop server
	// if err := srv.Stop(ctx); err != nil {
	//     log.Fatalf("Failed to stop archive server gracefully: %v", err)
	// }

	log.Println("Server stopped successfully")
}

// printHelp prints help information
func printHelp() {
	log.Printf(`%s %s - Data Persistence and Cache Invalidation Service for YaoVerse

Usage: %s [options]

Options:
  -config string
        Path to configuration file (default "%s")
  -version
        Print version information
  -help
        Print this help information

Environment Variables:
  YAO_ARCHIVE_CONFIG    Configuration file path (overrides -config flag)
  YAO_ARCHIVE_LOG_LEVEL Log level (debug, info, warn, error)

Examples:
  %s                           # Use default configuration
  %s -config /etc/archive.yaml # Use custom configuration file
  %s -version                  # Print version information

For more information, visit: https://github.com/eggybyte-technology/yao-verse
`, AppName, AppVersion, os.Args[0], DefaultConfigPath, os.Args[0], os.Args[0], os.Args[0])
}
