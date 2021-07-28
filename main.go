package main

import (
	"fmt"
	"log"
	"net/http"

	"nobl9-recruitment-task/calculator"
	"nobl9-recruitment-task/logging"
	"nobl9-recruitment-task/random"
	"nobl9-recruitment-task/rest"
)

const port = 8080

func main() {
	handler := rest.RandomHandler{Calc: calculator.NewCalculator(random.NewDao())}
	http.HandleFunc("/random/mean", handler.Handle)
	logging.Info("starting http server at port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
