// crawler crawls urls
package crawler

import (
	"net/url"
	"errors"
	"net/http"
	"golang.com/x/net/html"
	"github.com/fueledbymarvin/gocardless/logs"
)

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

func Crawl(u *url.URL) (string, error) {

	return "placeholder", nil
}
