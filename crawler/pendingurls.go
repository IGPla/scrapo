package crawler

import (
	"sync"
)

var pendingUrls []string
var pendingMutex = &sync.Mutex{}

// Add pending url
func addPendingURLs(urls []string) {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()
	for _, url := range urls {
		pendingUrls = append(pendingUrls, url)
	}
}

// Check if there are pending urls
func newURLsPending() bool {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()
	return len(pendingUrls) > 0
}

// Get pending urls pack
func getPendingURLsPack(packSize int) []string {
	pendingMutex.Lock()
	defer pendingMutex.Unlock()
	if packSize > len(pendingUrls) {
		urlsPack := pendingUrls
		pendingUrls = nil
		return urlsPack
	} else {
		urlsPack := pendingUrls[0 : packSize-1]
		pendingUrls = pendingUrls[packSize:]
		return urlsPack
	}
}
