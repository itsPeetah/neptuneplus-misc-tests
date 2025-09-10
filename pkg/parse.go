package main

import (
	"net/http"
	"os"
	"strconv"
)

const (
	defaultCount      = 2
	defaultMode       = ModeSequential
	defaultUpperBound = 20_000
)

func parseQuery(req *http.Request) (mode string, count int, upperBound int) {
	query := req.URL.Query()

	if paramMode := query.Get("mode"); paramMode == string(ModeSequential) || paramMode == string(ModeParallel) {
		mode = paramMode
	} else {
		mode = string(defaultMode)
	}

	if paramCount := query.Get("count"); len(paramCount) < 1 {
		count = defaultCount
	} else {
		_count, err := strconv.Atoi(paramCount)
		if err != nil {
			count = defaultCount
		} else {
			count = _count
		}
	}

	if paramUpperBound := query.Get("upperBound"); len(paramUpperBound) < 1 {
		upperBound = defaultUpperBound
	} else {
		_upper, err := strconv.Atoi(paramUpperBound)
		if err != nil {
			upperBound = defaultUpperBound
		} else {
			upperBound = _upper
		}
	}

	return mode, count, upperBound
}

func getBaseUri(varName string) (uri string, ok bool) {
	uri = os.Getenv(varName)
	return uri, len(uri) > 0
}
