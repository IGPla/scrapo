package crawler

import (
	"testing"
)

func TestAddVisitedUrl(t *testing.T) {
	var prevUrls int = len(visitedUrls)
	var url string = "test"
	addVisitedURL(url)
	var afterUrls int = len(visitedUrls)
	if afterUrls != prevUrls+1 {
		t.Errorf("Expected increment on visited urls. e(%d), g(%d)",
			prevUrls+1,
			afterUrls)
	}
	var lastUrl string = visitedUrls[len(visitedUrls)-1]
	if lastUrl != url {
		t.Errorf("Visited url do not match. e(%v), g(%v)",
			url,
			lastUrl)
	}
}

func TestIsUrlProcessed(t *testing.T) {
	// Reset visitedUrls
	visitedUrls = make([]string, 0)

	var url string = "test"
	if isUrlProcessed(url) {
		t.Errorf("Url should not be visited yet")
	}
	addVisitedURL(url)
	if !isUrlProcessed(url) {
		t.Errorf("Url should be marked as visited")
	}
}
