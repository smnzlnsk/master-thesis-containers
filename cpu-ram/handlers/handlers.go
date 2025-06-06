package handlers

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"benchmarking/benchmark"
)

// URL pattern for CPU activation with core count
var cpuActivatePattern = regexp.MustCompile(`^/cpu/activate(?:/(\d+))?$`)

// URL pattern for memory activation with memory limit
var memoryActivatePattern = regexp.MustCompile(`^/memory/activate(?:/(\d+))?$`)

// Version of the application - set from main
var BuildVersion = "0.0.1"

// HelloHandler responds with a simple greeting
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello, World!")
}

// HealthCheckHandler responds with a status message
func HealthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Server is up and running!")
}

// ActivateHandler handles CPU benchmark activation requests
// Supports both /activate and /cpu/activate[/cores] paths
func ActivateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the URL specifies a core count
	cores := 0 // Default to all cores

	// Check for the /cpu/activate/N pattern
	matches := cpuActivatePattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 && matches[1] != "" {
		// Extract the core count
		coreCount, err := strconv.Atoi(matches[1])
		if err != nil || coreCount <= 0 {
			http.Error(w, "Invalid core count", http.StatusBadRequest)
			return
		}
		cores = coreCount
	}

	if !benchmark.StartTaskWithCores(cores) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "CPU benchmark task is already running")
		return
	}

	// Give the task a moment to start and update its internal state
	time.Sleep(100 * time.Millisecond)

	// Now get the actual cores being used
	coresUsed := benchmark.GetCPUCoresUsed()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "CPU benchmark task activated successfully using %d cores", coresUsed)
}

// DeactivateHandler handles CPU benchmark deactivation requests
func DeactivateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !benchmark.StopTask() {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "No CPU benchmark task is currently running")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "CPU benchmark task deactivated successfully")
}

// ActivateMemoryHandler handles memory benchmark activation requests
// Supports /memory/activate[/limit] where limit is in MB (default: 1024 MB)
func ActivateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Check if the URL specifies a memory limit
	memoryLimit := 0 // Default to 1GB (set in the benchmark package)

	// Check for the /memory/activate/N pattern
	matches := memoryActivatePattern.FindStringSubmatch(r.URL.Path)
	if len(matches) > 1 && matches[1] != "" {
		// Extract the memory limit
		limit, err := strconv.Atoi(matches[1])
		if err != nil || limit <= 0 {
			http.Error(w, "Invalid memory limit", http.StatusBadRequest)
			return
		}
		memoryLimit = limit
	}

	if !benchmark.StartMemoryTaskWithLimit(memoryLimit) {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "Memory benchmark task is already running")
		return
	}

	// Get the actual memory limit being used
	limit := benchmark.GetMaxMemoryMB()
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Memory benchmark task activated successfully with %d MB limit", limit)
}

// DeactivateMemoryHandler handles memory benchmark deactivation requests
func DeactivateMemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if !benchmark.StopMemoryTask() {
		w.WriteHeader(http.StatusConflict)
		fmt.Fprintf(w, "No memory benchmark task is currently running")
		return
	}

	// Get the currently allocated memory
	allocatedMB := benchmark.GetAllocatedMemoryMB()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Memory benchmark task deactivated successfully. %d MB still allocated - use /memory/free to release.", allocatedMB)
}

// FreeMemoryHandler explicitly forces memory cleanup
func FreeMemoryHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get memory allocation before cleanup
	allocatedMB := benchmark.GetAllocatedMemoryMB()

	// Call the memory cleanup function
	benchmark.FreeAllMemory()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Forced memory cleanup completed. %d MB has been released back to the system.", allocatedMB)
}

// StatusHandler provides information about running benchmark tasks
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	cpuActive := benchmark.IsTaskRunning()
	memoryActive := benchmark.IsMemoryTaskRunning()

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Benchmark Status (Version: %s):\n", BuildVersion)
	fmt.Fprintf(w, "- CPU Benchmark: %s", statusText(cpuActive))

	if cpuActive {
		cores := benchmark.GetCPUCoresUsed()
		fmt.Fprintf(w, " (using %d cores)", cores)
	}
	fmt.Fprintf(w, "\n")

	// Always show memory info since memory can be allocated even when the task is not running
	allocatedMB := benchmark.GetAllocatedMemoryMB()
	fmt.Fprintf(w, "- Memory Benchmark: %s", statusText(memoryActive))

	if allocatedMB > 0 {
		maxMB := benchmark.GetMaxMemoryMB()
		percentage := 0
		if maxMB > 0 {
			percentage = allocatedMB * 100 / maxMB
		}

		if memoryActive {
			fmt.Fprintf(w, " (using %d MB of %d MB limit - %d%%)",
				allocatedMB, maxMB, percentage)
		} else {
			fmt.Fprintf(w, " (stopped, but still holding %d MB of memory)", allocatedMB)
		}
	}
	fmt.Fprintf(w, "\n")
}

// Helper function to convert boolean to status text
func statusText(active bool) string {
	if active {
		return "RUNNING"
	}
	return "STOPPED"
}
