package processor

import (
	"github.com/IGPla/scrapo/tasks"
	"testing"
)

func TestDownloadResourceProcessor(t *testing.T) {
	var rawHTML []byte = []byte("<html><head></head><body><div><a href='#test'>link</a><img src='test.jpg' /></div></body></html>")
	var task *tasks.Task = tasks.NewTask("http://test.com")
	task.Content = rawHTML
	var htmlProcessor *DownloadResourceProcessor = new(DownloadResourceProcessor)
	parsedResources, newRawResources := htmlProcessor.Parse(task)

	if len(parsedResources) != 1 {
		t.Errorf("Wrong returned parsed resources. e(%d), g(%d)",
			1,
			len(parsedResources))
	}
	if len(newRawResources) != 2 {
		t.Errorf("Wrong returned new raw resources. e(%d), g(%d)",
			1,
			len(newRawResources))
	}

	rawBytes1 := rawHTML
	rawBytes2 := parsedResources[0].Content
	if len(rawBytes1) != len(rawBytes2) {
		t.Errorf("DownloadResourceProcessor unexpected processed html size. e(%v), g(%v)", len(rawBytes1), len(rawBytes2))
	}
	for index, _ := range rawBytes1 {
		if rawBytes1[index] != rawBytes2[index] {
			t.Errorf("DownloadResourceProcessor unexpected change between raw html and processed html. e(%v), g(%v)", rawBytes1[index], rawBytes2[index])
		}
	}

}
