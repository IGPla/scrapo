package config

import (
	"encoding/json"
	"fmt"
	"github.com/IGPla/scrapo/logger"
	"github.com/IGPla/scrapo/processor"
	"github.com/IGPla/scrapo/storage"
	"io/ioutil"
	"log"
	"os"
)

var configLogger *log.Logger

func init() {
	configLogger = logger.GetLogger("CONFIG", os.Stdout)
}

/* Main configuration struct */
type ScrapoConfig struct {
	UserAgent         string                      `json:"user_agent"`
	URLs              []string                    `json:"urls"`
	Storage           storage.Storage             `json:"-"`
	AllowedDomains    []string                    `json:"allowed_domains"`
	NumWorkers        int                         `json:"num_workers"`
	ResourceProcessor processor.ResourceProcessor `json:"-"`
}

var MainConfig *ScrapoConfig

func PopulateMainConfig(configFilePath string) {
	MainConfig = getMainConfig(configFilePath)
}

func getMainConfig(configFilePath string) *ScrapoConfig {
	var mainConfig *ScrapoConfig = new(ScrapoConfig)

	jsonFile, openErr := os.Open(configFilePath)
	if openErr != nil {
		configLogger.Printf("Error arised trying to read config file %v (%v)",
			configFilePath,
			openErr.Error())
		panic("Could not parse configuration")
	}
	defer jsonFile.Close()

	byteJson, _ := ioutil.ReadAll(jsonFile)

	jsonErr := json.Unmarshal([]byte(byteJson), mainConfig)
	if jsonErr != nil {
		configLogger.Printf("Error arised trying to parse config file %v (%v)",
			configFilePath,
			jsonErr.Error())
		panic("Could not parse configuration")
	}

	return mainConfig
}

/* Utils test function to create a test config file */
func CreateTestConfigFile() (string, string, string, string, int, error) {
	var ua string = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/69.0.3497.100 Safari/537.36"
	var url string = "https://github.com/"
	var domain string = "github.com"
	var workers int = 10
	var content string = fmt.Sprintf("{\"user_agent\": \"%v\",\"urls\": [\"%v\"],\"allowed_domains\": [\"%v\"],\"num_workers\": %d}", ua, url, domain, workers)

	// Store content into file under a filepath
	var filepath string = "/tmp/test_config_scrapo.json"
	err := ioutil.WriteFile(filepath, []byte(content), 0700)
	if err != nil {
		return "", "", "", "", -1, err
	}
	return filepath, ua, url, domain, workers, nil
}
