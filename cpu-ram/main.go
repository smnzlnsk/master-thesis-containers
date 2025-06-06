package main

import (
	"fmt"
	"log"
	"net/http"

	"benchmarking/config"
	"benchmarking/handlers"
)

// buildVersion will be set during build via -ldflags
var buildVersion = "0.0.1" // Default version if not set during build

func main() {
	// Get application configuration
	cfg := config.GetDefaultConfig()
	serverAddr := fmt.Sprintf("%s:%s", cfg.ServerHost, cfg.ServerPort)

	// Pass version information to handlers package
	handlers.BuildVersion = buildVersion

	// Log version information
	log.Printf("Starting CPU-RAM benchmarking server version %s", buildVersion)

	// Register routes
	http.HandleFunc("/", handlers.HelloHandler)
	http.HandleFunc("/health", handlers.HealthCheckHandler)

	// CPU benchmark endpoints - using flexible pattern matching in the handler
	http.HandleFunc("/cpu/activate", handlers.ActivateHandler)
	http.HandleFunc("/cpu/activate/", handlers.ActivateHandler) // To handle /cpu/activate/N
	http.HandleFunc("/cpu/deactivate", handlers.DeactivateHandler)

	// Memory benchmark endpoints - using flexible pattern matching in the handler
	http.HandleFunc("/memory/activate", handlers.ActivateMemoryHandler)
	http.HandleFunc("/memory/activate/", handlers.ActivateMemoryHandler) // To handle /memory/activate/N
	http.HandleFunc("/memory/deactivate", handlers.DeactivateMemoryHandler)
	http.HandleFunc("/memory/free", handlers.FreeMemoryHandler) // Endpoint to explicitly free memory

	// Status endpoint
	http.HandleFunc("/status", handlers.StatusHandler)

	// Version endpoint to display container version
	http.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "CPU-RAM Benchmark Server Version: %s\n", buildVersion)
	})

	// Legacy endpoints (for backward compatibility)
	http.HandleFunc("/activate", handlers.ActivateHandler)
	http.HandleFunc("/deactivate", handlers.DeactivateHandler)

	// Start the server
	fmt.Printf("Server starting on %s...\n", serverAddr)
	if err := http.ListenAndServe(serverAddr, nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}
