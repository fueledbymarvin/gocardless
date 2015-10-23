// crawler crawls urls
package crawler

import (
	"net/url"
	"errors"
	"net/http"
	"golang.com/x/net/html"
	"github.com/fueledbymarvin/gocardless/logs"
)

type Sitemap struct {
	Host string
	Nodes map[*url.URL]*Node
}

type Node struct {
	URL *url.URL
	Neighbors []*Node
}

func Parse(uStr string) (*url.URL, error) {
	u, err := url.Parse(uStr)
	if err != nil {
		return nil, err
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return nil, errors.New("URL is not http or https.")
	}

	return u, nil
}

func Crawl(u *url.URL) ([]*Node, error) {

	sitemap := &Sitemap{Host: u.Host, Nodes: make(map[*url.URL]*Node)}
	sitemap.Nodes[u] = &Node{URL: u, Neighbors: make([]*Node, 0)}

	urls := make([]*url.URL, 0)
	urls = append(urls, u)
	for len(urls) > 0 {
		var toCrawl *url.URL
		toCrawl, urls = urls[0], urls[1:len(urls)]
		links := sitemap.process(toCrawl)
		urls = append(urls, links...)
	}
	
	return sitemap.asSlice(), nil
}

func (this *Sitemap) process(u *url.URL) []*url.URL {
	
	links := getLinks(u)

	// Only return unseen links with the given Host
	newLinks := make([]*url.URL, len(links))
	for _, link := range(links) {
		if _, seen := this.Nodes[link]; !seen && link.Host == this.Host {
			newLinks = append(newLinks, link)
		}
	}

	// Create links between original node and discovered links
	node :=	this.Nodes[u]
	for _, link := range(links) {
		var linkedNode *Node
		var seen bool
		if linkedNode, seen = this.Nodes[link]; !seen {
			linkedNode := &Node{URL: link, Neighbors: make([]*Node, 1)}
			this.Nodes[link] = linkedNode
		}
		linkedNode.Neighbors = append(linkedNode.Neighbors, node)
		node.Neighbors = append(node.Neighbors, linkedNode)
	}

	return newLinks
}

func (this *Sitemap) asSlice() []*Node {

	nodes := make([]*Node, 0)
	for _, node := range(this.Nodes) {
		nodes = append(nodes, node)
	}
	
	return nodes
}

func getLinks(u *url.URL) []*url.URL {

	return make([]*url.URL, 0)
}

