package requests

import (
	"github.com/IGPla/scrapo/config"
	"os"
	"testing"
)

func TestGetResource(t *testing.T) {
	// Prepare config
	filepath, _, _, _, _, error := config.CreateTestConfigFile()

	if error != nil {
		t.Errorf("Could not create the json config test file. (%v)",
			error.Error())
		return
	}
	config.PopulateMainConfig(filepath)
	// Wrong url
	var url string = "htt://test"
	buffer, _, _, _ := GetResource(url)
	if buffer != nil {
		t.Errorf("Expected error due to wrong url but got bytes. (%v)",
			buffer)
	}

	// Right url
	url = "https://www.google.com/"
	buffer, statusCode, _, error := GetResource(url)
	if error != nil {
		t.Errorf("Unexpected error arised when getting resource. (%v)",
			error.Error())
	}
	if statusCode != 200 {
		t.Errorf("Wrong status code received. e(%d), g(%d)",
			200,
			statusCode)
	}
	if len(buffer.Bytes()) <= 0 {
		t.Errorf("Expected response, empty line received")
	}
	os.Remove(filepath)
}
