package storage

import (
	"github.com/IGPla/scrapo/tasks"
	"os"
	"testing"
)

func TestFileStoreData(t *testing.T) {
	// Base content
	var content string = "Test content"
	var contentBytes = []byte(content)
	var task *tasks.Task = new(tasks.Task)
	task.URL = "http://www.google.com/test1/"
	task.Content = contentBytes
	task.ContentType = "text/html; charset=UTF-8"

	// Filestorage init
	var fs *FileStorage = new(FileStorage)
	fs.Prefix = "/tmp/testfilestorage"
	defer os.RemoveAll(fs.Prefix)

	// Store file without name
	htmlError := fs.StoreData(task)
	if htmlError != nil {
		t.Errorf("Errors appeared while storing content. (%v)",
			htmlError.Error())
	}
	var expectedFilepath string = "/tmp/testfilestorage/www.google.com/test1/data.htm"
	_, error := os.Stat(expectedFilepath)
	if error != nil {
		t.Errorf("File without name could not be stored. e(%v)",
			expectedFilepath)
	}

	// Store file with name
	task.URL = "http://www.google.com/test2?p1=v1"
	htmlError = fs.StoreData(task)
	if htmlError != nil {
		t.Errorf("Errors appeared while storing content. (%v)",
			htmlError.Error())
	}
	expectedFilepath = "/tmp/testfilestorage/www.google.com/test2?p1=v1"
	_, error = os.Stat(expectedFilepath)
	if error != nil {
		t.Errorf("File with name could not be stored. e(%v)",
			expectedFilepath)
	}

}
