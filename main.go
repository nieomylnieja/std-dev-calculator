package main

import (
	"log"
	"net/http"

	"nobl9-recruitment-task/calculator"
	"nobl9-recruitment-task/random"
	"nobl9-recruitment-task/rest"
)

func main() {
	handler := rest.RandomHandler{Calc: calculator.NewCalculator(random.NewDao())}
	http.HandleFunc("/random/mean", handler.Handle)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
