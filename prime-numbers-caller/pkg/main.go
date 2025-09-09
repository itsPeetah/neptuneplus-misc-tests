package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

var (
	baseUriSequential = ""
	baseUriParallel   = ""
)

func init() {
	seq, okSeq := getBaseUri("PRIME_NUMBERS_URL_SEQUENTIAL")
	par, okPar := getBaseUri("PRIME_NUMBERS_URL_PARALLEL")

	if !okSeq || !okPar {
		log.Fatal("prime-numbers URI not set")
	} else {
		baseUriSequential = seq
		baseUriParallel = par
	}

}

func main() {
	http.HandleFunc("/health", handleHealth)
	http.HandleFunc("/_/ready", handleReady)
	http.HandleFunc("/entrypoint", handlePrime)

	addr := fmt.Sprintf(":%d", 8080)
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
	mode, count, upperBound := parseQuery(r)

	resp := ""

	switch mode {
	case "seq":
		resp = callSequential(count, upperBound)
	case "par":
		resp = callParallel(count, upperBound)
	default:
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("mode parameter missing or invalid"))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(resp))
}
