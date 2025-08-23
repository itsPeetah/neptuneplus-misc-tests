package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
)

var (
	url  = ""
	port = 8080
)

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/_/ready", handleReady)
	http.HandleFunc("/entrypoint", handlePrime)

	if p, err := strconv.Atoi(os.Getenv("SERVICE_PORT")); err != nil {
		port = p
	}

	addr := fmt.Sprintf(":%d", port)
	log.Print("prime-numbers function starting on port 8080")
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("prime-numbers function failed to start: %v", err)
		os.Exit(2)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	log.Printf("received health check request")
	w.Write([]byte("health"))
}

func handleReady(w http.ResponseWriter, r *http.Request) {
	log.Printf("received ready check request")
	w.Write([]byte("ready"))
}

func handlePrime(w http.ResponseWriter, r *http.Request) {

	q := r.URL.Query()

	// count, err := strconv.Atoi(q.Get("count"))
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadRequest)
	// 	w.Write([]byte("Invalid count parameter."))
	// 	return
	// }

	count := 5

	switch q.Get("mode") {
	case "seq":
		url = os.Getenv("PRIME_NUMBERS_URL_SEQUENTIAL")

		if len(url) < 1 {
			log.Fatalf("No prime-numbers URL defined in env vars.")
			os.Exit(1)
		}

		callSequential(count)
	case "par":

		url = os.Getenv("PRIME_NUMBERS_URL_PARALLEL")

		if len(url) < 1 {
			log.Fatalf("No prime-numbers URL defined in env vars.")
			os.Exit(1)
		}

		callParallel(count)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing mode parameter."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(r.URL.RawQuery))
}

func callSequential(count int) {

	endpoint := fmt.Sprintf("%s/%d", url, 100000)

	log.Printf("Calling %s %d times sequentially", endpoint, count)
	for i := 0; i < count; i++ {
		doRequest(endpoint)
	}
}

func callParallel(count int) {
	log.Printf("Calling %s %d times in parallel", url, count)

	var wg sync.WaitGroup
	wg.Add(count)

	callFunc := func(multiplier int) {
		defer wg.Done()
		endpoint := fmt.Sprintf("%s/%d", url, 20000*multiplier)
		doRequest(endpoint)
	}

	for i := 0; i < count; i++ {
		go callFunc(i + 1)
	}

	wg.Wait()
}

func doRequest(endpoint string) {
	_, err := http.Get(endpoint)

	if err != nil {
		log.Printf("Error while making request to %s: %v", endpoint, err)
	}
}
