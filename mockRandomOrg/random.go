package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
)

// all of this because I was banned by some perf & security tool on random.org
func main() {
	http.HandleFunc("/integers", func(w http.ResponseWriter, r *http.Request) {
		num, err := readIntQueryParam(r, "num")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		min, err := readIntQueryParam(r, "min")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		max, err := readIntQueryParam(r, "max")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		for i := 0; i < num; i++ {
			_, _ = w.Write([]byte(strconv.Itoa(rand.Intn(max-min)+min) + "\n"))
		}
		w.WriteHeader(http.StatusOK)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
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
