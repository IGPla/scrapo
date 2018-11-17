package crawler

import (
	"github.com/IGPla/scrapo/config"
	"testing"
)

func TestSanitizeUrl(t *testing.T) {
	var raw string = "  http://test.com  "
	var expected string = "http://test.com"
	var got string = sanitizeUrl(raw)
	if got != expected {
		t.Errorf("Wrong transformation. e(%v), g(%v)",
			expected,
			got)
	}
}

func TestIsUrlValid(t *testing.T) {
	// Prepare config
	filepath, _, url, _, _, error := config.CreateTestConfigFile()

	if error != nil {
		t.Errorf("Could not create the json config test file. (%v)",
			error.Error())
		return
	}
	config.PopulateMainConfig(filepath)

	if !isUrlValid(url) {
		t.Errorf("Expected url to be valid. (%v)",
			url)
	}
	addVisitedURL(url)
	if isUrlValid(url) {
		t.Errorf("Expected url to be non valid. (%v)",
			url)
	}
	url = "htt://invalid"
	if isUrlValid(url) {
		t.Errorf("Expected url to be non valid. (%v)",
			url)
	}
}
