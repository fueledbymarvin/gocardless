# Simple Domain-Restricted Web Crawler

A simple, parallel web crawler written in Go with a web interface. It crawls a single domain and outputs a sitemap.

## Using the web crawler

1. Install [Go](https://golang.org/doc/install)
2. Download the package `go get github.com/fueledbymarvin/gocardless`
3. Navigate to the package `cd $GOPATH/src/github.com/fueledbymarvin/gocardless/`
4. Install dependencies `go get`
5. Start the server `go run server.go`
6. Go to [http://localhost:8000](http://localhost:8000)
7. Type in a URL and hit crawl
8. You can run tests with `go test`

## Design Considerations

I chose to implement the web crawler in Go due to the simple concurrency primitives. The initial, non-parallel version of the application took around six minutes to crawl [http://gocardless.com](http://gocardless.com) (~1000 pages). Most of the time is spent making http requests so there's a lot of opportunity for performance gains through concurrency. I used the main thread for link management (i.e. tracking which pages have been seen and updating the sitemap) and a worker pool to actually crawl the links that are discovered and placed in a jobs channel. The concurrent version takes around one minute to run on my machine (the number of workers is based on the number of logical CPUs of the server, so on my machine it was 8).

The crawler considers a, img, link, and script tags and looks for URLs in the href and src attributes.

The front-end displays the URLs in the order visited and lists all the links on that page. The request for the information is asynchronous since it takes a relatively long time and is displayed once the domain has been fully crawled. The logs show the progress of the request (e.g., how many URLs crawled, how many URLs pending). Given more time, I would create a more informative implementation that would stream the results as they occur (e.g., through websockets).
