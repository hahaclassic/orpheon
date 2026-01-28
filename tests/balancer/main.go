package main

import (
	"fmt"
	"net/http"
	"sort"
	"sync"
	"time"
)

type Result struct {
	Upstream string
	Error    error
}

type Stats struct {
	mu        sync.RWMutex
	upstreams map[string]int
	errors    int
	noHeader  int
	total     int
}

func main() {
	url := "http://localhost/api/v1/albums"
	numRequests := 100
	maxWorkers := 1

	fmt.Printf("n=%d requests to %s\n", numRequests, url)

	stats := &Stats{
		upstreams: make(map[string]int),
	}

	jobs := make(chan int, numRequests)
	results := make(chan Result, numRequests)

	wg := &sync.WaitGroup{}
	for w := 0; w < maxWorkers; w++ {
		wg.Add(1)
		go worker(w, url, jobs, results, wg)
	}

	startTime := time.Now()
	for i := 0; i < numRequests; i++ {
		jobs <- i
	}
	close(jobs)

	wgCollector := &sync.WaitGroup{}
	wgCollector.Add(1)

	go func() {
		defer wgCollector.Done()
		for result := range results {
			stats.mu.Lock()
			stats.total++

			if result.Error != nil {
				stats.errors++
			} else if result.Upstream == "" {
				stats.noHeader++
			} else {
				stats.upstreams[result.Upstream]++
			}
			stats.mu.Unlock()
		}
	}()

	wg.Wait()
	close(results)
	wgCollector.Wait()

	endTime := time.Now()
	duration := endTime.Sub(startTime)

	fmt.Printf("Time: %.2f sec\n", duration.Seconds())
	fmt.Println("Stat:")
	fmt.Println("==========================================")

	stats.mu.RLock()
	defer stats.mu.RUnlock()

	type kv struct {
		key   string
		value int
	}

	var sorted []kv
	for k, v := range stats.upstreams {
		sorted = append(sorted, kv{k, v})
	}

	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].value > sorted[j].value
	})

	servers := []string{"content-msv", "content-msv-ro-1", "content-msv-ro-2"}

	for i, item := range sorted {
		percentage := float64(item.value) / float64(numRequests) * 100
		fmt.Printf("%-15s (%-16s): %3d запросов (%5.1f%%)\n",
			item.key, servers[i], item.value, percentage)
	}

	fmt.Println("==========================================")
	fmt.Printf("Total requests: %d\n", stats.total)
	fmt.Printf("Errors: %d\n", stats.errors)
	fmt.Printf("No header: %d\n", stats.noHeader)
}

func worker(id int, url string, jobs <-chan int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()

	for range jobs {
		client := &http.Client{
			Timeout: 10 * time.Second,
		}

		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			results <- Result{Error: err}
			continue
		}

		req.Header.Add("Connection", "close")
		req.Header.Add("Cache-Control", "no-cache")
		req.Header.Add("Pragma", "no-cache")
		req.Header.Add("User-Agent", fmt.Sprintf("Analyzer/%d", id))

		resp, err := client.Do(req)
		if err != nil {
			results <- Result{Error: err}
			continue
		}

		upstream := resp.Header.Get("X-Debug-Upstream")
		resp.Body.Close()

		results <- Result{Upstream: upstream}
	}
}
