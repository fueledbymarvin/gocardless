// crawler crawls urls
package crawler

import (
	"errors"
	"github.com/fueledbymarvin/gocardless/logs"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"fmt"
	"runtime"
)

const chunkSize int = 100

type Sitemap struct {
	Host  string
	Nodes map[string]*Node
}

type Node struct {
	URL       *url.URL
	Neighbors map[string]bool
}

type result struct {
	URL   *url.URL
	Links []*url.URL
}

func Parse(uStr string) (*url.URL, error) {

	u, err := url.Parse(uStr)
	if err != nil {
		return nil, err
	}

	if !ensureCanonical(u) {
		return nil, errors.New("Protocol is not http or https.")
	}
	return u, nil
}

func ensureCanonical(u *url.URL) bool {

	if u.Scheme != "http" && u.Scheme != "https" {
		return false
	}
	// consider https as http to reduce redundancy
	u.Scheme = "http"

	// remove trailing slash to reduce redundancy
	if len(u.Path) != 0 && u.Path[len(u.Path)-1] == '/' {
		u.Path = u.Path[:len(u.Path)-1]
	}

	// clear fragment to reduce redundancy
	u.Fragment = ""

	return true
}

func Crawl(u *url.URL) (*Sitemap, error) {

	sitemap := &Sitemap{Host: u.Host, Nodes: make(map[string]*Node)}
	sitemap.Nodes[u.String()] = &Node{URL: u, Neighbors: make(map[string]bool)}

	// create incoming and outgoing channels for workers
	urls := make(chan *url.URL, 1000)
	results := make(chan *result, 100)

	// create worker pool
	nWorkers := runtime.NumCPU()
	runtime.GOMAXPROCS(nWorkers)
	for i := 0; i < nWorkers; i++ {
		go worker(urls, results)
	}

	// add initial url and track outstanding jobs to know when to terminate
	urls <- u
	outstanding := 1

	for count := 1; ; count++ {
		res := <-results
		newLinks := sitemap.update(res)
		outstanding += len(newLinks) - 1
		for _, link := range newLinks {
			urls <- link
		}
		
		if outstanding == 0 {
			close(urls)
			logs.Log(fmt.Sprintf("Crawled %d urls total", count))
			break
		}

		if count % chunkSize == 0 {
			logs.Log(fmt.Sprintf("Crawled %d urls so far", count))
			logs.Log(fmt.Sprintf("%d urls pending", outstanding))
		}
	}

	return sitemap, nil
}

func worker(urls <-chan *url.URL, results chan<- *result) {

	for u := range urls {
		results <- &result{URL: u, Links: getLinks(u)}
	}
}

func (this *Sitemap) update(res *result) []*url.URL {

	// Only return unseen links with the given Host
	node := this.Nodes[res.URL.String()]
	newLinks := make([]*url.URL, 0, len(res.Links))
	for _, link := range res.Links {
		var linkedNode *Node
		var seen bool
		if linkedNode, seen = this.Nodes[link.String()]; !seen {
			linkedNode = &Node{URL: link, Neighbors: make(map[string]bool)}
			this.Nodes[link.String()] = linkedNode

			if link.Host == this.Host {
				newLinks = append(newLinks, link)
			}
		}
		node.Neighbors[linkedNode.URL.String()] = true
	}

	return newLinks
}

func getLinks(u *url.URL) []*url.URL {

	resp, err := http.Get(u.String())
	if err != nil {
		logs.Log(fmt.Sprintf("Couldn't crawl %s", u))
	}
	defer resp.Body.Close()

	links := make([]*url.URL, 0)
	tokenizer := html.NewTokenizer(resp.Body)
	for {
		tokenType := tokenizer.Next()
		switch tokenType {
		case html.ErrorToken:
			return links
		case html.StartTagToken, html.SelfClosingTagToken:
			token := tokenizer.Token()
			if link, ok := getURL(u, token); ok {
				links = append(links, link)
			}
		}
	}

	return links
}

func getURL(src *url.URL, token html.Token) (*url.URL, bool) {

	if !contains([]string{"a", "img", "script", "link"}, token.Data) {
		return nil, false
	}

	for _, attr := range token.Attr {
		if contains([]string{"href", "src"}, attr.Key) {
			if u, err := src.Parse(attr.Val); err == nil && ensureCanonical(u) {
				return u, true
			}
		}
	}

	return nil, false
}

func contains(s []string, str string) bool {
	
	for _, a := range s {
		if a == str {
			return true
		}
	}
	return false
}
