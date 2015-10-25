// Starts the server on port 8000.
package main

import (
	"github.com/fueledbymarvin/gocardless/handlers"
	"github.com/fueledbymarvin/gocardless/logs"
	"net/http"
)

func init() {
	logs.Initialize("gocardless")
	http.HandleFunc("/crawl", handlers.BeforeAction(handlers.Crawl))
}

func main() {
	logs.Log("Starting server on port 8000")
	logs.CheckFatal(http.ListenAndServe(":8000", nil))
}
