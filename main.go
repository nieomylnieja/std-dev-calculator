package main

import (
	"fmt"
	"log"
	"net/http"

	"nobl9-recruitment-task/calculator"
	"nobl9-recruitment-task/random"
	"nobl9-recruitment-task/rest"
)

const port = 8080

func main() {
	handler := rest.RandomHandler{Calc: calculator.NewCalculator(random.NewDao())}
	http.HandleFunc("/random/mean", handler.Handle)
	log.Printf("starting http server at port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
