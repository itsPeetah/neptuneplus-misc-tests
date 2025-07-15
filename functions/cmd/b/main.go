package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

func handler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Function B received a request from %s", r.RemoteAddr)
	response := "Hello from Function B!\n"

	urls := make([]string, 2)
	urls[0] = os.Getenv("RW_URL")
	urls[1] = os.Getenv("W2_URL")

	var wg sync.WaitGroup
	wg.Add(2)

	callFunction := func(url string, parallel bool) {

		if parallel {
			defer wg.Done()
		}

		if url != "" {
			log.Printf("Function B calling %s", url)
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
			log.Println("Url not set. Function B is the end of this chain.")
			response += "Url not set. Function B is the end of this chain.\n"
		}
	}

	go callFunction(urls[0], true)
	go callFunction(urls[1], true)

	wg.Wait()

	callFunction(urls[1], false)

	fmt.Fprint(w, response)
}

func main() {
	http.HandleFunc("/", handler)
	port := "8080"
	log.Printf("Function B server starting on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Function B server failed to start: %v", err)
	}
}
