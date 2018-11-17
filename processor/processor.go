package processor

import (
	"bytes"
	"github.com/IGPla/scrapo/logger"
	"github.com/IGPla/scrapo/tasks"
	"log"
	"os"
)

var processorLogger *log.Logger

func init() {
	processorLogger = logger.GetLogger("PROCESSOR", os.Stdout)
}

// This will be the base class to extend crawler. Each client should implement these interfaces to add functionality to each type
type ResourceProcessor interface {
	Parse(task *tasks.Task) (parsedResources []*tasks.Task, newRawResources []*tasks.Task)
}

// Minimal implementation of ResourceProcessor
type DownloadResourceProcessor struct{}

func (dhp DownloadResourceProcessor) Parse(task *tasks.Task) (parsedResources []*tasks.Task, newRawResources []*tasks.Task) {
	var links []string = GetHTMLElements(bytes.NewBuffer(task.Content), "a", "href")
	var images []string = GetHTMLElements(bytes.NewBuffer(task.Content), "img", "src")
	var scripts []string = GetHTMLElements(bytes.NewBuffer(task.Content), "script", "src")
	var csss []string = GetHTMLElements(bytes.NewBuffer(task.Content), "link", "href")

	processorLogger.Printf("Got %d links, %d images, %d scripts and %d css",
		len(links), len(images), len(scripts), len(csss))

	newRawResources = make([]*tasks.Task, 0, len(images)+len(links)+len(scripts)+len(csss))

	for _, link := range links {
		newRawResources = append(newRawResources, tasks.NewTask(link))
	}
	for _, image := range images {
		newRawResources = append(newRawResources, tasks.NewTask(image))
	}
	for _, script := range scripts {
		newRawResources = append(newRawResources, tasks.NewTask(script))
	}
	for _, css := range csss {
		newRawResources = append(newRawResources, tasks.NewTask(css))
	}
	parsedResources = make([]*tasks.Task, 0, 1)
	parsedResources = append(parsedResources, task)
	return parsedResources, newRawResources
}
