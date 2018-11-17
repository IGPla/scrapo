package main

import (
	"github.com/IGPla/scrapo/config"
	"github.com/IGPla/scrapo/crawler"
	"github.com/IGPla/scrapo/processor"
	"github.com/IGPla/scrapo/storage"
)

func main() {
	// Base settings
	config.PopulateMainConfig("sampleconfig.json")
	// Storage choice and configuration
	var fileStorage *storage.FileStorage = new(storage.FileStorage)
	fileStorage.Prefix = "/tmp/myscrapingproject"
	config.MainConfig.Storage = fileStorage
	// Base HTML class to parse content
	config.MainConfig.ResourceProcessor = new(processor.DownloadResourceProcessor)
	// Start crawler
	crawler.Start()
}
