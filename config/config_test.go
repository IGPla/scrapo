package config

import (
	"os"
	"testing"
)

func TestPopulateMainConfig(t *testing.T) {
	filepath, ua, url, domain, workers, error := CreateTestConfigFile()
	if error != nil {
		t.Errorf("Could not create the json config test file. (%v)",
			error.Error())
		return
	}
	// Use that filepath to call PopulateMainConfig
	PopulateMainConfig(filepath)

	// Access to MainConfig var and check each property
	if MainConfig.NumWorkers != workers {
		t.Errorf("Number of workers do not match. e(%d), g(%d)",
			MainConfig.NumWorkers,
			workers)
	}
	if len(MainConfig.URLs) != 1 || MainConfig.URLs[0] != url {
		t.Errorf("URLs do not match. e(%v), g(%v)",
			url,
			MainConfig.URLs)
	}
	if len(MainConfig.AllowedDomains) != 1 || MainConfig.AllowedDomains[0] != domain {
		t.Errorf("Domains do not match. e(%v), g(%v)",
			domain,
			MainConfig.AllowedDomains)
	}
	if MainConfig.UserAgent != ua {
		t.Errorf("User agent do not match. e(%v), g(%v)",
			ua,
			MainConfig.UserAgent)
	}

	// Delete created file
	os.Remove(filepath)
}
