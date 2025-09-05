package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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

	count := 2
	resp := ""

	switch q.Get("mode") {
	case "seq":
		url = os.Getenv("PRIME_NUMBERS_URL_SEQUENTIAL")
		url = strings.TrimSuffix(url, "/")

		if len(url) < 1 {
			log.Fatalf("No prime-numbers URL defined in env vars.")
			os.Exit(1)
		}

		resp = callSequential(count)
	case "par":

		url = os.Getenv("PRIME_NUMBERS_URL_PARALLEL")
		url = strings.TrimSuffix(url, "/")

		if len(url) < 1 {
			log.Fatalf("No prime-numbers URL defined in env vars.")
			os.Exit(1)
		}

		resp = callParallel(count)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Missing mode parameter."))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}

func callSequential(count int) string {

	endpoint := url
	response := ""

	log.Printf("Calling %s %d times sequentially", endpoint, count)
	for i := 0; i < count; i++ {
		response += doRequest(endpoint) + "\n"
	}
	return response
}

func callParallel(count int) string {
	log.Printf("Calling %s %d times in parallel", url, count)
	response := ""
	var wg sync.WaitGroup
	wg.Add(count)

	callFunc := func(multiplier int) {
		defer wg.Done()
		endpoint := url
		response += doRequest(endpoint) + "\n"
	}

	for i := 0; i < count; i++ {
		go callFunc(i + 1)
	}

	wg.Wait()

	return response
}

func doRequest(endpoint string) string {
	_, err := http.Get(endpoint)

	if err != nil {
		r := fmt.Sprintf("Error while making request to %s: %v\n", endpoint, err)
		log.Print(r)
	}
	return "OK"
}
