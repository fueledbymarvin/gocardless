package crawler_test

import (
	"github.com/fueledbymarvin/gocardless/logs"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCrawler(t *testing.T) {
	logs.Initialize("gocardless (test)")

	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawler Suite")
}
