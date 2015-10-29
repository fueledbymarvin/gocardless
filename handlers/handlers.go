// handlers provides handler functions for endpoints that return json data
package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/fueledbymarvin/gocardless/crawler"
	"github.com/fueledbymarvin/gocardless/logs"
	"net/http"
	"net/http/httptest"
)

// BeforeAction logs request and response data and times the handler
// h's execution.
func BeforeAction(h func(http.ResponseWriter, *http.Request), contentType string) func(http.ResponseWriter, *http.Request) {

	return func(rw http.ResponseWriter, req *http.Request) {
		defer logs.TimerEnd(logs.TimerBegin(fmt.Sprintf("%s '%s'", req.Method, req.URL.Path)))

		// log request headers
		logs.Log(fmt.Sprintf("Request headers: %s", req.Header))

		// parse params
		err := req.ParseForm()
		if logs.CheckErr(err) {
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		// new recorder for logging/middleware
		rec := httptest.NewRecorder()
		// set content-type
		if contentType != "" {
			rec.Header().Set("Content-Type", contentType)
		}
		
		h(rec, req)

		// log response
		logs.Log(fmt.Sprintf("Response status: %d", rec.Code))
		logs.Log(fmt.Sprintf("Response headers: %s", rec.Header()))

		// copy to actual ResponseWriter
		copyResponse(rw, rec)
	}
}

// copyResponse copies all relevant info from rec to rw.
func copyResponse(rw http.ResponseWriter, rec *httptest.ResponseRecorder) {

	// copy the headers
	for k, v := range rec.Header() {
		rw.Header()[k] = v
	}
	// copy the code
	rw.WriteHeader(rec.Code)
	// copy the body
	rw.Write(rec.Body.Bytes())
}

// JSON marshal's the response variable into json and prints it on rw.
func JSON(rw http.ResponseWriter, response interface{}) {

	encoded, err := json.Marshal(response)
	if logs.CheckErr(err) {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, string(encoded))
}

// GET /crawl
// Crawls a given url
func Crawl(rw http.ResponseWriter, req *http.Request) {

	uStr := req.Form.Get("url")
	if uStr == "" {
		http.Error(rw, "Missing url parameter", 422)
		return
	}

	u, err := crawler.Parse(uStr)
	if logs.CheckErr(err) {
		http.Error(rw, err.Error(), 422)
		return
	}

	JSON(rw, crawler.Crawl(u))
}

func Index(rw http.ResponseWriter, req *http.Request) {

	http.ServeFile(rw, req, "assets/views/index.html")
}

func Assets(rw http.ResponseWriter, req *http.Request) {

	http.ServeFile(rw, req,
		fmt.Sprintf("assets/%s/%s", req.Form.Get("type"), req.Form.Get("file")))
}
