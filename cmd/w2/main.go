package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Function W2 received a request from %s", r.RemoteAddr)
	response := "Hello from Function W2!\n"

	delayDuration := time.Duration(2) * time.Second
	time.Sleep(delayDuration)

	response += "Slept for 2 seconds.\n"

	fmt.Fprint(w, response)
}

func main() {
	http.HandleFunc("/", handler)
	port := "8080"
	log.Printf("Function W2 server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Function W2 server failed to start: %v", err)
	}
}
