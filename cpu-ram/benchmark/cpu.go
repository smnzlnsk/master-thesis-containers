package benchmark

import (
	"fmt"
	"math"
	"math/rand"
	"runtime"
	"sync"
	"time"
)

// Global variables to control the CPU benchmark task
var (
	cpuTaskRunning bool
	cpuTaskChan    chan bool
	cpuTaskWg      sync.WaitGroup
	cpuTaskMutex   sync.Mutex
	numCoresUsed   int // Number of CPU cores currently being used
)

// init initializes the package-level variables
func init() {
	cpuTaskRunning = false
	cpuTaskChan = make(chan bool, 1) // Buffered channel to prevent blocking
	numCoresUsed = 0
	rand.Seed(time.Now().UnixNano())
}

// performCPUIntensiveMath does CPU-intensive calculations to generate load infinitely
func performCPUIntensiveMath(stopChan chan bool) float64 {
	result := 0.0
	// Generate a random base number
	base := rand.Float64() * 100

	// Counter to periodically check for stop signal
	counter := 0

	// Run indefinitely until signaled to stop
	for {
		// Check if we need to stop more frequently
		select {
		case <-stopChan:
			return result
		default:
			// Continue processing
		}

		// Perform expensive math operations
		x := base + math.Sin(float64(counter)/1000)
		result += math.Sin(x) * math.Cos(x) * math.Exp(math.Sin(x/5))
		result += math.Sqrt(math.Abs(result)) + math.Pow(x, 0.5*math.Sin(x))

		// Reset result occasionally to prevent overflow
		if math.Abs(result) > 1e10 {
			result = rand.Float64() * 100
		}

		// Check if we need to stop periodically
		counter++
		if counter%1000 == 0 { // Check more frequently (was 10000)
			select {
			case <-stopChan:
				return result
			default:
				// Continue processing
			}
		}
	}
}

// startCPUTask runs CPU-intensive calculations continuously until signaled to stop
func startCPUTask(coreCount int) {
	defer cpuTaskWg.Done()

	// If coreCount is invalid or zero, use all cores
	if coreCount <= 0 {
		coreCount = runtime.NumCPU()
	}

	// Make sure we don't exceed the number of available cores
	availableCores := runtime.NumCPU()
	if coreCount > availableCores {
		coreCount = availableCores
	}

	// Store the number of cores in use - use mutex to ensure visibility across goroutines
	cpuTaskMutex.Lock()
	numCoresUsed = coreCount
	cpuTaskMutex.Unlock()

	fmt.Printf("CPU benchmark task started - generating load using %d of %d available CPU cores\n",
		coreCount, availableCores)

	// Create a ticker for status updates
	statusTicker := time.NewTicker(10 * time.Second)
	defer statusTicker.Stop()

	// Create channel for results
	resultChan := make(chan float64, coreCount) // Make this buffered

	// Function for worker goroutines
	worker := func(id int, stopChan chan bool) {
		for {
			select {
			case <-stopChan:
				fmt.Printf("Worker %d stopping\n", id)
				return
			default:
				result := performCPUIntensiveMath(stopChan)
				// Send result but don't block if no one is listening
				select {
				case resultChan <- result:
				default:
				}
			}
		}
	}

	// Start worker goroutines (limited by coreCount)
	workerStopChans := make([]chan bool, coreCount)
	for i := 0; i < coreCount; i++ {
		workerStopChans[i] = make(chan bool, 1) // Buffered channel for each worker
		go worker(i, workerStopChans[i])
	}

	var totalCalcs uint64

	// Main control loop
	for {
		select {
		case <-cpuTaskChan:
			// Signal all workers to stop
			fmt.Println("CPU benchmark task received stop signal, shutting down all workers...")

			for i, stopChan := range workerStopChans {
				fmt.Printf("Stopping worker %d...\n", i)
				stopChan <- true
				close(stopChan)
			}

			fmt.Printf("CPU benchmark task stopped after %d calculation cycles\n", totalCalcs)
			return

		case result := <-resultChan:
			// Just keep track of calculations and occasionally use the result
			// to prevent the compiler from optimizing away the work
			totalCalcs++
			if totalCalcs%1000 == 0 {
				// Use the result to prevent optimization
				fmt.Printf("\rPerformed %d calculation cycles (last result: %.5g)...", totalCalcs, result)
			}

		case <-statusTicker.C:
			fmt.Printf("\nCPU benchmark running - using %d cores - completed %d calculation cycles so far\n",
				coreCount, totalCalcs)
		}
	}
}

// StartTaskWithCores starts the CPU benchmark task using the specified number of cores
// Returns true if task was started, false if it was already running
func StartTaskWithCores(cores int) bool {
	cpuTaskMutex.Lock()
	defer cpuTaskMutex.Unlock()

	if cpuTaskRunning {
		return false
	}

	// Create a fresh channel for this task
	cpuTaskChan = make(chan bool, 1)

	// Set numCoresUsed based on the requested cores
	if cores <= 0 {
		numCoresUsed = runtime.NumCPU() // Default to all cores
	} else {
		availableCores := runtime.NumCPU()
		if cores > availableCores {
			numCoresUsed = availableCores
		} else {
			numCoresUsed = cores
		}
	}

	// Start the CPU task with specified core count
	cpuTaskRunning = true
	cpuTaskWg.Add(1)
	go startCPUTask(cores)

	return true
}

// StartTask starts the CPU benchmark task using all available cores
// Returns true if task was started, false if it was already running
func StartTask() bool {
	return StartTaskWithCores(0) // 0 means use all available cores
}

// StopTask stops the CPU benchmark task if it's running
// Returns true if task was stopped, false if it wasn't running
func StopTask() bool {
	cpuTaskMutex.Lock()

	if !cpuTaskRunning {
		cpuTaskMutex.Unlock()
		return false
	}

	fmt.Println("Sending stop signal to CPU benchmark task...")

	// Signal the task to stop
	select {
	case cpuTaskChan <- true:
		fmt.Println("Stop signal sent successfully")
	default:
		fmt.Println("Warning: Channel was full, but proceeding with shutdown")
	}

	cpuTaskRunning = false

	// Reset cores used count
	numCoresUsed = 0

	// Unlock before waiting to avoid deadlock
	cpuTaskMutex.Unlock()

	fmt.Println("Waiting for CPU task to complete shutdown...")
	cpuTaskWg.Wait()
	fmt.Println("CPU task shutdown complete")

	return true
}

// IsTaskRunning returns the current state of the CPU benchmark task
func IsTaskRunning() bool {
	cpuTaskMutex.Lock()
	defer cpuTaskMutex.Unlock()
	return cpuTaskRunning
}

// GetCPUCoresUsed returns the number of CPU cores currently being used by the benchmark
func GetCPUCoresUsed() int {
	cpuTaskMutex.Lock()
	defer cpuTaskMutex.Unlock()
	return numCoresUsed
}
