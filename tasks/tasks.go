package tasks

import ()

type Task struct {
	URL         string
	ContentType string
	Content     []byte
}

func NewTask(url string) *Task {
	var newTask *Task = new(Task)
	newTask.URL = url
	return newTask
}
