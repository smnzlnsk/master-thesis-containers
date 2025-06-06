package benchmark

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

// Global variables to control the memory benchmark task
var (
	memoryTaskRunning bool
	memoryTaskChan    chan bool
	memoryTaskWg      sync.WaitGroup
	memoryTaskMutex   sync.Mutex

	// Memory storage
	memoryBlocks      [][]byte
	memoryBlocksMutex sync.Mutex
	maxMemoryMB       int // Maximum memory to allocate in MB
)

// Memory allocation sizes
const (
	blockSize          = 10 * 1024 * 1024 // 10MB per block
	allocationDelay    = 500 * time.Millisecond
	statusInterval     = 5 * time.Second
	defaultMaxMemoryMB = 1024 // Default max memory is 1GB (1024MB)
)

// init initializes the package-level variables
func init() {
	memoryTaskRunning = false
	memoryTaskChan = make(chan bool, 1) // Use buffered channel to prevent blocking
	maxMemoryMB = defaultMaxMemoryMB
}

// SetMaxMemoryMB sets the maximum amount of memory in MB that the benchmark can allocate
// Returns the actual limit that was set
func SetMaxMemoryMB(mbLimit int) int {
	memoryTaskMutex.Lock()
	defer memoryTaskMutex.Unlock()

	if mbLimit <= 0 {
		maxMemoryMB = defaultMaxMemoryMB
	} else {
		maxMemoryMB = mbLimit
	}

	return maxMemoryMB
}

// GetMaxMemoryMB returns the current maximum memory limit in MB
func GetMaxMemoryMB() int {
	memoryTaskMutex.Lock()
	defer memoryTaskMutex.Unlock()
	return maxMemoryMB
}

// freeMemory explicitly releases memory by clearing the memory blocks
// and forcing garbage collection
func freeMemory() {
	memoryBlocksMutex.Lock()
	// Get current memory allocation for reporting
	allocatedMB := 0
	if memoryBlocks != nil {
		allocatedMB = len(memoryBlocks) * blockSize / (1024 * 1024)
	}

	// Explicitly set to nil to release references
	memoryBlocks = nil
	memoryBlocksMutex.Unlock()

	// Force garbage collection
	fmt.Printf("Cleaning up %d MB of allocated memory...\n", allocatedMB)
	runtime.GC()

	// Wait a moment and force another GC for good measure
	time.Sleep(500 * time.Millisecond)
	runtime.GC()
	fmt.Println("Memory cleanup complete - memory should now be released to the system")
}

// startMemoryTask continuously allocates memory until signaled to stop or reaching the limit
func startMemoryTask() {
	defer memoryTaskWg.Done()

	// Get the memory limit at task start time
	memoryTaskMutex.Lock()
	memoryLimit := maxMemoryMB
	memoryTaskMutex.Unlock()

	fmt.Printf("Memory benchmark task started - will allocate up to %d MB\n", memoryLimit)

	// Initialize the memory blocks slice if it doesn't exist
	memoryBlocksMutex.Lock()
	if memoryBlocks == nil {
		// Pre-allocate capacity for the slice
		memoryBlocks = make([][]byte, 0, memoryLimit/(blockSize/(1024*1024)))
	}
	// Calculate current allocation
	currentAllocation := 0
	if len(memoryBlocks) > 0 {
		currentAllocation = len(memoryBlocks) * blockSize / (1024 * 1024)
		fmt.Printf("Reusing existing memory allocation of %d MB\n", currentAllocation)
	}
	memoryBlocksMutex.Unlock()

	// Create a ticker for memory allocation and status updates
	allocTicker := time.NewTicker(allocationDelay)
	statusTicker := time.NewTicker(statusInterval)
	defer allocTicker.Stop()
	defer statusTicker.Stop()

	allocatedMB := currentAllocation

	// Main control loop
	for {
		select {
		case <-memoryTaskChan:
			// Stop the task but do NOT free memory
			fmt.Printf("Memory benchmark task stopped after allocating %d MB\n", allocatedMB)
			fmt.Println("Memory is still allocated. Use /memory/free endpoint to release it.")
			return

		case <-allocTicker.C:
			// Check if we've reached the memory limit
			if allocatedMB >= memoryLimit {
				fmt.Printf("\nReached memory allocation limit of %d MB. Stopping further allocations.\n", memoryLimit)
				// Keep the task running, but stop allocating more memory
				allocTicker.Stop()
				continue
			}

			// Allocate a new memory block
			memoryBlocksMutex.Lock()
			// Create a new memory block and fill it with data to ensure it's actually allocated
			newBlock := make([]byte, blockSize)
			for i := 0; i < len(newBlock); i += 1024 { // Fill every 1KB to ensure allocation
				newBlock[i] = byte(i % 256)
			}
			memoryBlocks = append(memoryBlocks, newBlock)
			currentBlocks := len(memoryBlocks)
			allocatedMB = currentBlocks * blockSize / (1024 * 1024)
			memoryBlocksMutex.Unlock()

			fmt.Printf("\rAllocated %d MB of memory (%d%% of limit)...",
				allocatedMB, allocatedMB*100/memoryLimit)

		case <-statusTicker.C:
			// Get current memory usage
			memoryBlocksMutex.Lock()
			currentBlocks := len(memoryBlocks)
			allocatedMB = currentBlocks * blockSize / (1024 * 1024)
			memoryBlocksMutex.Unlock()

			fmt.Printf("\nMemory benchmark running - using approximately %d MB (%d%% of %d MB limit)\n",
				allocatedMB, allocatedMB*100/memoryLimit, memoryLimit)
		}
	}
}

// StartMemoryTask starts the memory-intensive benchmark task
// Returns true if task was started, false if it was already running
func StartMemoryTask() bool {
	return StartMemoryTaskWithLimit(defaultMaxMemoryMB)
}

// StartMemoryTaskWithLimit starts the memory-intensive benchmark task with a specific MB limit
// If limit is <= 0, the default limit (1024 MB) is used
// Returns true if task was started, false if it was already running
func StartMemoryTaskWithLimit(mbLimit int) bool {
	memoryTaskMutex.Lock()
	defer memoryTaskMutex.Unlock()

	if memoryTaskRunning {
		return false
	}

	// Set the memory limit
	if mbLimit > 0 {
		maxMemoryMB = mbLimit
	} else {
		maxMemoryMB = defaultMaxMemoryMB
	}

	// Create a new channel for this task
	memoryTaskChan = make(chan bool, 1)

	// Start the memory task
	memoryTaskRunning = true
	memoryTaskWg.Add(1)
	go startMemoryTask()

	return true
}

// StopMemoryTask stops the memory-intensive benchmark task without freeing memory
// Returns true if task was stopped, false if it wasn't running
func StopMemoryTask() bool {
	memoryTaskMutex.Lock()

	if !memoryTaskRunning {
		memoryTaskMutex.Unlock()
		return false
	}

	// Signal the task to stop
	select {
	case memoryTaskChan <- true:
		fmt.Println("Stop signal sent to memory task")
	default:
		fmt.Println("Warning: Channel was full, but proceeding with shutdown")
	}

	memoryTaskRunning = false

	// Unlock before waiting to avoid deadlock
	memoryTaskMutex.Unlock()

	fmt.Println("Waiting for memory task to complete shutdown...")
	memoryTaskWg.Wait()
	fmt.Println("Memory task shutdown complete - memory is still allocated")

	return true
}

// IsMemoryTaskRunning returns the current state of the memory benchmark task
func IsMemoryTaskRunning() bool {
	memoryTaskMutex.Lock()
	defer memoryTaskMutex.Unlock()
	return memoryTaskRunning
}

// GetAllocatedMemoryMB returns the current amount of memory allocated by the benchmark in MB
func GetAllocatedMemoryMB() int {
	memoryBlocksMutex.Lock()
	defer memoryBlocksMutex.Unlock()

	if memoryBlocks == nil {
		return 0
	}

	return len(memoryBlocks) * blockSize / (1024 * 1024)
}

// FreeAllMemory is a public function that can be called to explicitly free memory
// even outside the normal benchmark stop flow
func FreeAllMemory() {
	freeMemory()
}
