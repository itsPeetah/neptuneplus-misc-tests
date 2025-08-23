package main

import (
	"fmt"
	"net/http"
)

const (
	CALLS_TO_MAKE = 10000
	ENDPOINT      = "http://localhost:8080/function/openfaas-fn/prime-numbers-caller-sequential/entrypoint?mode=seq"
)

func main() {

	for i := 0; i < CALLS_TO_MAKE; i++ {
		fmt.Printf("Call no. %d\n", i+1)
		_, err := http.Get(ENDPOINT)
		if err != nil {
			fmt.Println("error")
		}
	}
}
