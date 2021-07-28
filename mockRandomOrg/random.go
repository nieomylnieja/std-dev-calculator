package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"nobl9-recruitment-task/logging"
)

const port = 8081

// all of this because I was banned by some perf & security tool on random.org
func main() {
	sleepFor := getSleep()
	http.HandleFunc("/integers", func(w http.ResponseWriter, r *http.Request) {
		logging.Info("handling request from: %s", r.RequestURI)
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
		time.Sleep(sleepFor)
		for i := 0; i < num; i++ {
			_, _ = w.Write([]byte(strconv.Itoa(rand.Intn(max-min)+min) + "\n"))
		}
	})
	logging.Info("starting http mock-random server at port %d with sleep: %s", port, sleepFor.String())
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
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

func getSleep() time.Duration {
	sleepStr := os.Getenv("RANDOM_ORG_MOCK_SLEEP")
	if sleep, err := time.ParseDuration(sleepStr); err != nil {
		return 0
	} else {
		return sleep
	}
}
