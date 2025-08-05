package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Function RW received a request from %s", r.RemoteAddr)
	response := "Hello from Function RW!\n"

	delaySeconds := rand.Intn(10) + 1
	delayDuration := time.Duration(delaySeconds) * time.Second
	time.Sleep(delayDuration)

	response += fmt.Sprintf("Slept for %d seconds.\n", delaySeconds)

	fmt.Fprint(w, response)
}

func main() {
	http.HandleFunc("/handle", handler)
	http.HandleFunc("/_/ready", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "OK")
	})
	port := "8080"
	log.Printf("Function RW server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Function RW server failed to start: %v", err)
	}
}
