// Package main is the entry point for the GophKeeper CLI application.
// It initializes the database connection, configures and starts the HTTP server,
// and provides a CLI interface for versioning information.
//
// Features:
// - Starts an HTTP server to handle API requests.
// - Provides version and build date information via CLI.
// - Handles graceful shutdown on interrupt signals.
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gleb-korostelev/GophKeeper/cmd/initConnection"
	"github.com/gleb-korostelev/GophKeeper/config"
	"github.com/gleb-korostelev/GophKeeper/tools/closer"
	"github.com/gleb-korostelev/GophKeeper/tools/logger"
	"github.com/spf13/cobra"
)

var (
	version   = "dev"     // Application version, injected at build time.
	buildDate = "unknown" // Build date, injected at build time.
)

// main initializes the CLI application and HTTP server.
// It supports a "version" command to display the application version and build date.
// The function also handles server startup and graceful shutdown on interrupt signals.
func main() {
	// Defer the cleanup of all resources using the closer utility.
	defer func() {
		closer.Wait()
		closer.CloseAll()
	}()

	// Root command for the CLI application.
	rootCmd := &cobra.Command{
		Use:   "gophkeeper-cli-app",
		Short: "CLI for bank card application",
	}

	// Command to display version and build date.
	versionCmd := &cobra.Command{
		Use:   "version",
		Short: "Show version and build date",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", version)
			fmt.Printf("Build Date: %s\n", buildDate)
		},
	}

	// Add the version command to the root command.
	rootCmd.AddCommand(versionCmd)

	// Execute the root command.
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}

	// Create a background context for managing the server lifecycle.
	ctx := context.Background()

	// Initialize the database connection using the initConnection package.
	db := initConnection.NewDBConn(ctx)

	// Retrieve the server port from the configuration.
	port := config.GetConfigInt(config.Port)

	// Log the startup message, including the server port and Swagger documentation port.
	logger.Info(fmt.Sprintf("Started server at :%d. Swagger docs stated at %d", port, port+1))

	// Configure the HTTP server.
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: initConnection.InitImpl(ctx, db, port),
	}

	// Channel to capture OS interrupt signals.
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, os.Kill)

	// Start the HTTP server in a separate goroutine.
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal("Error starting server: %v", err)
		}
	}()

	// Block until an interrupt signal is received.
	<-stop

	// Gracefully shut down the server with a 5-second timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logger.Info("Shutting down server...")
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exiting")
}
