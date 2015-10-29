package crawler_test

import (
	. "github.com/fueledbymarvin/gocardless/crawler"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/gorilla/mux"
	"fmt"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Crawler", func() {

	Describe("Parse", func() {

		It("parses a url", func() {
			u, err := Parse("http://gocardless.com")
			Expect(err).To(BeNil())
			Expect(u.Host).To(Equal("gocardless.com"))
		})

		It("treats https as http", func() {
			u, err := Parse("https://gocardless.com")
			Expect(err).To(BeNil())
			Expect(u.Scheme).To(Equal("http"))
		})

		It("removes trailing slashes", func() {
			u, err := Parse("https://gocardless.com/blog/")
			Expect(err).To(BeNil())
			Expect(u.Path).To(Equal("/blog"))
		})

		It("removes fragments", func() {
			u, err := Parse("https://gocardless.com/blog#yolo")
			Expect(err).To(BeNil())
			Expect(u.Fragment).To(Equal(""))
		})

		It("errors if not http or https", func() {
			_, err := Parse("ws://gocardless.com")
			Expect(err).NotTo(BeNil())
		})

		It("errors if invalid url", func() {
			_, err := Parse("derp")
			Expect(err).NotTo(BeNil())
		})
	})

	Describe("Crawl", func() {

		makeNode := func(u string, offsite bool, ls ...string) map[string]interface{} {

			node := make(map[string]interface{})
			node["url"] = u
			node["offsite"] = offsite
			links := make([]string, 0, len(ls))
			links = append(links, ls...)
			node["links"] = links

			return node
		}

		var r *mux.Router
		var ts *httptest.Server

		JustBeforeEach(func() {
			r = mux.NewRouter()
			ts = httptest.NewServer(r)
		})

		It("crawls a url", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"/blog\">Blog</a>")
			})
			r.HandleFunc("/blog", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"/blog/yo\">yo</a>" +
					"<a href=\"/blog/sup\">sup</a>")
			})
			r.HandleFunc("/blog/yo", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "Yo!")
			})
			r.HandleFunc("/blog/sup", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "Sup?")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				ts.URL + "/blog",
			))
			expected = append(expected, makeNode(ts.URL + "/blog", false,
				ts.URL + "/blog/yo",
				ts.URL + "/blog/sup",
			))
			expected = append(expected, makeNode(ts.URL + "/blog/yo", false))
			expected = append(expected, makeNode(ts.URL + "/blog/sup", false))
			Expect(sitemap).To(Equal(expected))
		})

		It("handles img tags", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<img src=\"http://i.imgur.com/mDuAK1v.gif\" />")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				"http://i.imgur.com/mDuAK1v.gif",
			))
			expected = append(expected, makeNode("http://i.imgur.com/mDuAK1v.gif", true))
			Expect(sitemap).To(Equal(expected))
		})

		It("handles script tags", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<script src=\"/whatever.js\"></script>")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				ts.URL + "/whatever.js",
			))
			expected = append(expected, makeNode(ts.URL + "/whatever.js", false))
			Expect(sitemap).To(Equal(expected))
		})

		It("handles link tags", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<link href=\"/whatever.css\" />")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				ts.URL + "/whatever.css",
			))
			expected = append(expected, makeNode(ts.URL + "/whatever.css", false))
			Expect(sitemap).To(Equal(expected))
		})

		It("doesn't crawl offsite links", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"http://google.com\"></a>")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				"http://google.com",
			))
			expected = append(expected, makeNode("http://google.com", true))
			Expect(sitemap).To(Equal(expected))
		})

		It("doesn't fail when encountering http errors", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"/fail\"></a>")
			})
			r.HandleFunc("/fail", func(rw http.ResponseWriter, req *http.Request) {
				rw.WriteHeader(http.StatusInternalServerError)
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				ts.URL + "/fail",
			))
			expected = append(expected, makeNode(ts.URL + "/fail", false))
			Expect(sitemap).To(Equal(expected))
		})

		It("doesn't fail when encountering non-http links", func() {
			r.HandleFunc("/", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"/blog\">Blog</a>")
			})
			r.HandleFunc("/blog", func(rw http.ResponseWriter, req *http.Request) {
				fmt.Fprintln(rw, "<a href=\"ws://socket.io\">yo</a>")
			})

			u, _ := Parse(ts.URL)
			sitemap := Crawl(u)

			expected := make([]map[string]interface{}, 0)
			expected = append(expected, makeNode(ts.URL, false,
				ts.URL + "/blog",
			))
			expected = append(expected, makeNode(ts.URL + "/blog", false))
			Expect(sitemap).To(Equal(expected))
		})
	})
})
