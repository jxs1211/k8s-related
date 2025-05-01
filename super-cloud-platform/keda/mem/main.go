package main

import (
	"fmt"
	"net/http"
	"runtime"
	"time"
)

var memBuffer [][]byte

func main() {
	http.HandleFunc("/grow", func(w http.ResponseWriter, r *http.Request) {
		memBuffer = append(memBuffer, make([]byte, 10*1024*1024*1024))
		printMemStats()
		fmt.Fprintf(w, "Memory increased - Current: %d MB\n", len(memBuffer)*10)
	})

	http.HandleFunc("/shrink", func(w http.ResponseWriter, r *http.Request) {
		// Create new slice and force GC to release old memory
		// newBuffer := make([][]byte, len(memBuffer)/2)
		// copy(newBuffer, memBuffer[:len(memBuffer)/2])
		memBuffer = nil

		// Force garbage collection
		runtime.GC()
		time.Sleep(100 * time.Millisecond) // Give GC time to work

		printMemStats()
		fmt.Fprintf(w, "Memory decreased - Current: %d MB\n", len(memBuffer)*10)
		// write the mem status to the w

	})

	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		printMemStats()
		fmt.Fprintf(w, "Memory status - Current: %d MB\n", len(memBuffer)*10)
	})

	fmt.Println("Starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

func printMemStats() {
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Alloc = %v MiB", bToMb(m.Alloc))
	fmt.Printf("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
	fmt.Printf("\tSys = %v MiB", bToMb(m.Sys))
	fmt.Printf("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
	return b / 1024 / 1024
}
