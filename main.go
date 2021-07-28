package main

import (
	"fmt"
	"log"
	"net/http"

	"std-dev-calculator/calculator"
	"std-dev-calculator/logging"
	"std-dev-calculator/random"
	"std-dev-calculator/rest"
)

const port = 8080

func main() {
	handler := rest.RandomHandler{Calc: calculator.NewCalculator(random.NewDao())}
	http.HandleFunc("/random/mean", handler.Handle)
	logging.Info("starting http server at port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
