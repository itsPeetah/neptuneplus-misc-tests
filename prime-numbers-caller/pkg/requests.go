package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"sync"
	"time"
)

func doRequest(endpoint string) string {
	resp, err := http.Get(endpoint)

	if err != nil {
		r := fmt.Sprintf("Error while making request to %s: %v\n", endpoint, err)
		log.Print(r)
		return r
	}

	if resp.StatusCode == http.StatusOK {
		defer resp.Body.Close()
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Print(err)
			return fmt.Sprintf("Got reponse code %d but could not read response body", resp.StatusCode)
		}
		bodyString := string(bodyBytes)
		return bodyString
	} else {
		return fmt.Sprintf("Reponse was %d", resp.StatusCode)
	}
}

func buildEndpoint(baseUri string, upperBound int) (string, error) {
	endpoint, err := url.JoinPath(baseUri, fmt.Sprintf("/prime/%d", upperBound))
	return endpoint, err
}

func callSequential(count int, upperBound int) string {

	endpoint, err := buildEndpoint(baseUriSequential, upperBound)
	if err != nil {
		return fmt.Sprintf("Could not build the url: %v", err)
	}

	log.Printf("Calling %s %d times sequentially", endpoint, count)

	t0 := time.Now()

	response := ""
	for i := 0; i < count; i++ {
		r := doRequest(endpoint)
		response += r + "\n"
	}

	log.Printf("Finished in: %dms", time.Now().UnixMilli()-t0.UnixMilli())

	return response
}

func callParallel(count int, upperBound int) string {

	endpoint, err := buildEndpoint(baseUriParallel, upperBound)
	if err != nil {
		return fmt.Sprintf("Could not build the url: %v", err)
	}

	response := ""
	var wg sync.WaitGroup
	wg.Add(count)

	callFunc := func(id int) {
		defer wg.Done()
		response += fmt.Sprintf("%d - %s\n", id, doRequest(endpoint))
	}

	log.Printf("Calling %s %d times in parallel", endpoint, count)

	t0 := time.Now()

	for i := 0; i < count; i++ {
		go callFunc(i + 1)
	}

	wg.Wait()

	log.Printf("Finished in: %dms", time.Now().UnixMilli()-t0.UnixMilli())

	return response
}
