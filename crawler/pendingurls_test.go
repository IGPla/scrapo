package crawler

import (
	"testing"
)

func TestAddPendingUrl(t *testing.T) {
	var prevUrls int = len(pendingUrls)
	var url string = "test"
	addPendingURLs([]string{url})
	var afterUrls int = len(pendingUrls)
	if afterUrls != prevUrls+1 {
		t.Errorf("Expected increment on pending urls. e(%d), g(%d)",
			prevUrls+1,
			afterUrls)
	}
	var lastUrl string = pendingUrls[len(pendingUrls)-1]
	if lastUrl != url {
		t.Errorf("Pending url do not match. e(%v), g(%v)",
			url,
			lastUrl)
	}
}

func TestNewUrlsPending(t *testing.T) {
	pendingUrls = make([]string, 0)
	if newURLsPending() {
		t.Errorf("It should not be pending urls")
	}
	addPendingURLs([]string{"test"})
	if !newURLsPending() {
		t.Errorf("It should show pending urls")
	}
}

func TestGetPendingURLsPack(t *testing.T) {
	pendingUrls = make([]string, 0)
	var pack []string = getPendingURLsPack(10)
	if len(pack) != 0 {
		t.Errorf("Expected empty pack")
	}
	addPendingURLs([]string{"test1", "test2"})
	pack = getPendingURLsPack(10)
	if len(pack) != 2 {
		t.Errorf("Wrong pack retrieved. e(%d), g(%d)",
			2,
			len(pack))
	}
	if newURLsPending() {
		t.Errorf("Expected empty pending url queue")
	}
}
