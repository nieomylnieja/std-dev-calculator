package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"std-dev-calculator/calculator"
	"std-dev-calculator/logging"
)

type stdDevCalculator interface {
	CalculateStdDev(ctx context.Context, requests, length int) ([]calculator.Result, error)
}

type RandomHandler struct {
	Calc stdDevCalculator
}

func (h RandomHandler) Handle(w http.ResponseWriter, r *http.Request) {
	reqBegin := time.Now()
	requests, err := readIntQueryParam(r, "requests")
	if err != nil {
		logging.Error("error occurred while reading 'requests' query param {err=%s}", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	length, err := readIntQueryParam(r, "length")
	if err != nil {
		logging.Error("error occurred while reading 'length' query param {err=%s}", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	logging.Info("handling GET /random/mean request {requests=%d, length=%d}", requests, length)

	res, err := h.Calc.CalculateStdDev(context.Background(), requests, length)
	if err != nil {
		// question: should our end api user see the underlying cause of the problem?
		// I've went with yes: but I was also wondering whether not to log it and return 500 on all occurrences of httpError
		if hErr, ok := err.(*httpError); ok {
			logging.Error(hErr.Error())
			http.Error(w, hErr.Message, hErr.StatusCode)
			return
		}
		// handle context timeout, we can rely on the timeout error message as it's part of the std lib (for now at least...)
		if ctxErr, ok := err.(*url.Error); ok && ctxErr.Err.Error() == "context deadline exceeded" {
			logging.Error("%s {url=%s}", ctxErr.Error(), ctxErr.URL)
			http.Error(w, "request has timed out", http.StatusServiceUnavailable)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		logging.Error("error occurred while encoding  query param {err=%s}", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	logging.Info("successfully handled GET /random/mean request in: %s", time.Now().Sub(reqBegin).String())
}

func readIntQueryParam(r *http.Request, name string) (int, error) {
	query := r.URL.Query()
	if _, ok := query[name]; !ok {
		return 0, fmt.Errorf("query param: %s is missing", name)
	}
	i, err := strconv.Atoi(query.Get(name))
	if err != nil {
		return 0, fmt.Errorf("%s query param must be a valid integer", name)
	}
	return i, nil
}
