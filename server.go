// Starts the server on port 8000.
package main

import (
	"github.com/fueledbymarvin/gocardless/handlers"
	"github.com/fueledbymarvin/gocardless/logs"
	"github.com/gorilla/mux"
	"net/http"
)

func init() {
	logs.Initialize("gocardless")
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/",
		handlers.BeforeAction(handlers.Index, "text/html")).Methods("GET")
	r.HandleFunc("/crawl",
		handlers.BeforeAction(handlers.Crawl, "application/json")).Methods("GET")
	r.PathPrefix("/assets/").Handler(
		http.StripPrefix("/assets/", http.FileServer(http.Dir("./assets/"))))

	logs.Log("Starting server on port 8000")
	logs.CheckFatal(http.ListenAndServe(":8000", r))
}
