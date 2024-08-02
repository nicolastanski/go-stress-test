package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

type Result struct {
	StatusCode int
	Duration   time.Duration
}

func worker(url string, requests int, results chan<- Result, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := 0; i < requests; i++ {
		start := time.Now()
		resp, err := http.Get(url)
		duration := time.Since(start)
		if err != nil {
			results <- Result{StatusCode: 0, Duration: duration}
			continue
		}
		results <- Result{StatusCode: resp.StatusCode, Duration: duration}
		resp.Body.Close()
	}
}

func main() {
	url := flag.String("url", "", "URL")
	totalRequests := flag.Int("requests", 0, "Total de Requisições")
	concurrency := flag.Int("concurrency", 0, "Chamadas concorrentes")

	flag.Parse()

	if *url == "" || *totalRequests <= 0 || *concurrency <= 0 {
		fmt.Println("Parâmetros inválidos.")
		os.Exit(1)
	}

	results := make(chan Result, *totalRequests)
	var wg sync.WaitGroup

	start := time.Now()
	for i := 0; i < *concurrency; i++ {
		requestsPerWorker := *totalRequests / *concurrency
		if i == *concurrency-1 {
			requestsPerWorker += *totalRequests % *concurrency
		}
		wg.Add(1)
		go worker(*url, requestsPerWorker, results, &wg)
	}

	wg.Wait()
	close(results)
	totalDuration := time.Since(start)

	statusCount := make(map[int]int)
	var successCount int
	for result := range results {
		if result.StatusCode == 200 {
			successCount++
		}
		statusCount[result.StatusCode]++
	}

	fmt.Printf("Tempo de Extecução: %v\n", totalDuration)
	fmt.Printf("Total de requisições: %d\n", *totalRequests)
	fmt.Printf("Total Requisições Status 200: %d\n", successCount)
	fmt.Println("Total Requisições Outros Status: ")

	for status, count := range statusCount {
		if status != 200 {
			fmt.Printf("Status %d: %d\n", status, count)
		}
	}
}
