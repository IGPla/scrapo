package crawler

import (
	"github.com/IGPla/scrapo/config"
	"testing"
)

func TestInitializeResources(t *testing.T) {
	// Prepare config
	filepath, _, _, _, _, error := config.CreateTestConfigFile()

	if error != nil {
		t.Errorf("Could not create the json config test file. (%v)",
			error.Error())
		return
	}
	config.PopulateMainConfig(filepath)

	// Check resources
	initializeResources()
	if len(pendingTasks) != 1 {
		t.Errorf("Expected 1 task pending in pending tasks. e(%d), g(%d)",
			1,
			len(pendingTasks))
	}
}

func TestUpdateWorkerStatus(t *testing.T) {
	workersStatus = make([]bool, 2)
	workersStatus[0] = false
	workersStatus[1] = false
	updateWorkerStatus(0, true)
	if !workersStatus[0] {
		t.Errorf("Expected worker status to change. e(%v), g(%v)",
			true,
			workersStatus[0])
	}
}

func TestCheckWorkersFinished(t *testing.T) {
	workersStatus = make([]bool, 2)
	workersStatus[0] = false
	workersStatus[1] = false
	if !checkWorkersFinished() {
		t.Errorf("Expected finished workers")
	}
	workersStatus[1] = true
	if checkWorkersFinished() {
		t.Errorf("Expected busy workers")
	}
}

func TestUpdatePendingTasksCounter(t *testing.T) {
	pendingTasksCounter = 0
	updatePendingTasksCounter(true)
	if pendingTasksCounter != 1 {
		t.Errorf("pendingTasksCounter did not update value as expected. e(%d), g(%d)",
			1,
			pendingTasksCounter)
	}
	updatePendingTasksCounter(false)
	if pendingTasksCounter != 0 {
		t.Errorf("pendingTasksCounter did not update value as expected. e(%d), g(%d)",
			0,
			pendingTasksCounter)
	}
}

func TestTaskCounterExhausted(t *testing.T) {
	pendingTasksCounter = 0
	if !taskCounterExhausted() {
		t.Errorf("Expected true from tasksCounterExhausted")
	}
	pendingTasksCounter = 5
	if taskCounterExhausted() {
		t.Errorf("Expected false from tasksCounterExhausted")
	}

}
