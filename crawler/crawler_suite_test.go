package crawler_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
	"github.com/fueledbymarvin/gocardless/logs"
)

func TestCrawler(t *testing.T) {
	logs.Initialize("gocardless (test)")
	
	RegisterFailHandler(Fail)
	RunSpecs(t, "Crawler Suite")
}
