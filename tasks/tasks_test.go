package tasks

import (
	"testing"
)

func TestNewTask(t *testing.T) {
	var testUrl string = "http://test.com"
	var task *Task = NewTask(testUrl)
	if task.URL != testUrl {
		t.Errorf("Wrong URL value returned by NewTask constructor. e(%v), g(%v)",
			testUrl,
			task.URL)
	}
}
