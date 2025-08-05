package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Function A received a request from %s", r.RemoteAddr)
	response := "Hello from Function A!\n"

	urls := make([]string, 2)
	urls[0] = os.Getenv("W2_URL")
	urls[1] = os.Getenv("RW_URL")

	for _, url := range urls {
		// Check for the next service URL (Function B)
		if url != "" {
			log.Printf("Function A calling %s", url)
			resp, err := http.Get(url)
			if err != nil {
				log.Printf("Error calling Function %s: %v", url, err)
				response += fmt.Sprintf("Error calling Function %s: %v\n", url, err)
			} else {
				defer resp.Body.Close()
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading response from %s: %v", url, err)
					response += fmt.Sprintf("Error reading response from %s: %v\n", url, err)
				} else {
					response += fmt.Sprintf("Response from %s:\n%s", url, string(body))
				}
			}
		} else {
			log.Println("Url not set. Function A is the end of this chain.")
			response += "Url not set. Function A is the end of this chain.\n"
		}
	}

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
	log.Printf("Function A server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Function A server failed to start: %v", err)
	}
}
