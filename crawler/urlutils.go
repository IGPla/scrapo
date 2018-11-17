package crawler

import (
	"github.com/IGPla/scrapo/config"
	"net/url"
	"strings"
)

// Sanitize url
func sanitizeUrl(url string) string {
	var _url string = strings.TrimSpace(url)
	return _url
}

/* Return true if url is valid, false otherwhise
Validity rules:
- Not visited yet
- In allowed domains
*/
func isUrlValid(_url string) bool {
	if !isUrlProcessed(_url) {
		parsedUrl, parseError := url.Parse(_url)
		if parseError != nil {
			crawlerLogger.Printf("Droping url by parse error %v (%v)",
				_url,
				parseError.Error())
			return false
		}
		for _, domain := range config.MainConfig.AllowedDomains {
			if parsedUrl.Host == domain {
				return true
			}
		}
	}
	return false
}
