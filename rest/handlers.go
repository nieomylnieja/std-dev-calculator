package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"nobl9-recruitment-task/calculator"
)

type stdDevCalculator interface {
	CalculateStdDev(ctx context.Context, requests, length int) ([]calculator.Result, error)
}

type RandomHandler struct {
	Calc stdDevCalculator
}

func (h RandomHandler) Handle(w http.ResponseWriter, r *http.Request) {
	requests, err := readIntQueryParam(r, "requests")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	length, err := readIntQueryParam(r, "length")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	res, err := h.Calc.CalculateStdDev(context.Background(), requests, length)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err = json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
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
