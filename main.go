package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

type Result struct {
	TotalRunTime  time.Duration
	TotalRequests int32
}

func main() {
	url, requests, concurrency := parseFlags()
	var wg sync.WaitGroup
	var mutex sync.Mutex
	var httpStatusDistribution = make(map[int]int)
	var countRequests atomic.Int32

	startTime := time.Now()
	for i := 0; i < requests; i++ {
		for j := 0; j < concurrency; j++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				defer func() {
					countRequests.Add(1)
				}()
				resp, err := http.Get(url)
				if err != nil {
					fmt.Println("Error fetching ", url, ":", err)
				}
				mutex.Lock()
				httpStatusDistribution[resp.StatusCode]++
				mutex.Unlock()
			}()
		}
	}
	wg.Wait()
	result := Result{
		TotalRunTime:  time.Since(startTime),
		TotalRequests: countRequests.Load(),
	}

	println("Total runtime:", result.TotalRunTime.String())
	println("Total requests:", result.TotalRequests)
	for statusCode, count := range httpStatusDistribution {
		println("Status Code:", statusCode, "Count:", count)
	}
}

func parseFlags() (string, int, int) {
	url := flag.String("url", "", "URL to send requests")
	requests := flag.Int("requests", 0, "Number of requests to send")
	concurrency := flag.Int("concurrency", 1, "Number of concurrent requests to send")
	flag.Parse()

	if *url == "" {
		fmt.Println("Error: --url is required")
		flag.Usage()
		os.Exit(1)
	}

	if *requests == 0 {
		fmt.Println("Error: --requests is required")
		flag.Usage()
		os.Exit(1)
	}

	if *concurrency == 0 {
		fmt.Println("Error: --concurrency is required")
		flag.Usage()
		os.Exit(1)
	}

	return *url, *requests, *concurrency
}
