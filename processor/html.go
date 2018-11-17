package processor

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
)

// Return all links inside resource
func GetHTMLElements(resource *bytes.Buffer, element string, attr string) []string {
	document, err := goquery.NewDocumentFromReader(resource)
	if err != nil {
		processorLogger.Printf("Error arised parsing links (%v)",
			err.Error())
		return nil
	}
	return processElements(document.Find(element), attr)
}

// Link processor for a goquery selection
func processElements(elements *goquery.Selection, attr string) []string {
	var results []string = make([]string, 0, len(elements.Nodes))
	elements.Each(func(index int, element *goquery.Selection) {
		attrVal, exists := element.Attr(attr)
		if exists {
			results = append(results, attrVal)
		}
	})
	return results
}
