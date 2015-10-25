// crawler crawls urls
package crawler

import (
	"errors"
	"github.com/fueledbymarvin/gocardless/logs"
	"golang.org/x/net/html"
	"net/http"
	"net/url"
	"fmt"
)

type Sitemap struct {
	Host  string
	Nodes map[string]*Node
}

type Node struct {
	URL       *url.URL
	Neighbors map[string]bool
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

	return true
}

func Crawl(u *url.URL) (*Sitemap, error) {

	sitemap := &Sitemap{Host: u.Host, Nodes: make(map[string]*Node)}
	sitemap.Nodes[u.String()] = &Node{URL: u, Neighbors: make(map[string]bool)}

	urls := make([]*url.URL, 0)
	urls = append(urls, u)
	for len(urls) > 0 {
		fmt.Println(len(urls))
		var toCrawl *url.URL
		toCrawl, urls = urls[0], urls[1:len(urls)]
		links := getLinks(toCrawl)
		fmt.Println(toCrawl)
		newLinks := sitemap.update(toCrawl, links)
		urls = append(urls, newLinks...)
	}

	return sitemap, nil
}

func (this *Sitemap) update(u *url.URL, links []*url.URL) []*url.URL {

	// Only return unseen links with the given Host
	node := this.Nodes[u.String()]
	newLinks := make([]*url.URL, 0)
	for _, link := range links {
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
