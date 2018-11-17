package crawler

import (
	"sync"
)

var visitedUrls []string
var visitedMutex = &sync.Mutex{}

// Add visited url
func addVisitedURL(url string) {
	visitedMutex.Lock()
	defer visitedMutex.Unlock()
	visitedUrls = append(visitedUrls, url)
}

// Check if url is already visited
func isUrlProcessed(url string) bool {
	visitedMutex.Lock()
	defer visitedMutex.Unlock()
	for _, _url := range visitedUrls {
		if _url == url {
			return true
		}
	}
	return false
}
