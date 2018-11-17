package crawler

import (
	"fmt"
	"github.com/IGPla/scrapo/config"
	"github.com/IGPla/scrapo/logger"
	"github.com/IGPla/scrapo/requests"
	"github.com/IGPla/scrapo/tasks"
	"log"
	"os"
	"sync"
	"time"
)

var workersStatus []bool
var pendingTasks chan *tasks.Task
var workerStatusMutex = &sync.Mutex{}
var pendingTasksCounterMutex = &sync.Mutex{}
var finishedTasks chan interface{}
var crawlerLogger *log.Logger
var pendingTasksCounter int

func init() {
	crawlerLogger = logger.GetLogger("CRAWLER", os.Stdout)
}

func crawlerLog(crawlerId int, message string) {
	crawlerLogger.Printf("Crawler %d: %v", crawlerId, message)
}

// Start crawler run
func Start() {
	initializeResources()
	initializeWorkers()
	<-finishedTasks
}

// Initialize resources
func initializeResources() {
	crawlerLogger.Printf("Initializing resources...")
	var urls []string = config.MainConfig.URLs
	visitedUrls = make([]string, len(urls))
	pendingUrls = make([]string, len(urls))
	finishedTasks = make(chan interface{})
	pendingTasks = make(chan *tasks.Task, 10*config.MainConfig.NumWorkers)
	pendingTasksCounter = len(urls)
	for _, url := range urls {
		pendingTasks <- tasks.NewTask(url)
	}
	crawlerLogger.Printf("Resources initialized")
}

// Initialize workers
func initializeWorkers() {
	crawlerLogger.Printf("Initializing workers...")
	var numWorkers int = config.MainConfig.NumWorkers
	workersStatus = make([]bool, numWorkers)
	for i := 0; i < numWorkers; i++ {
		go crawl(pendingTasks, i)
	}
	go crawlManager(pendingTasks, finishedTasks)
	crawlerLogger.Printf("Workers initialized")
}

// Crawl worker function
func crawl(pendingTasks chan *tasks.Task, crawlerId int) {
	updateWorkerStatus(crawlerId, false)
	for {
		select {
		case task, ok := <-pendingTasks:
			if ok {
				updateWorkerStatus(crawlerId, true)
				crawlerLog(crawlerId, fmt.Sprintf("Check workers finished: %v",
					checkWorkersFinished()))
				crawlerLog(crawlerId, fmt.Sprintf("Got new task: %v", task.URL))
				// Get resource
				resource, _, headers, requestError := requests.GetResource(task.URL)

				if requestError == nil {
					// Assign data
					task.ContentType = headers.Get("Content-Type")
					task.Content = resource.Bytes()

					// Apply processor
					parsedResources, newRawResources := config.MainConfig.ResourceProcessor.Parse(task)

					// Process new raw resources
					var newLinkBuffer []string = make([]string, 0, len(newRawResources))
					for _, nrr := range newRawResources {
						// Validate links
						if isUrlValid(nrr.URL) {
							// Store link as visited to prevent duplicated entries
							addVisitedURL(nrr.URL)
							// Add new link to pending urls
							newLinkBuffer = append(newLinkBuffer, nrr.URL)
						} else {
							crawlerLog(crawlerId, fmt.Sprintf("Droping link %v",
								nrr.URL))
						}
					}

					// Add all new links to pending urls
					addPendingURLs(newLinkBuffer)

					// Store resources
					for _, parsedTask := range parsedResources {
						config.MainConfig.Storage.StoreData(parsedTask)
					}
				}
				crawlerLog(crawlerId, fmt.Sprintf("Finished task: %v", task.URL))
				updateWorkerStatus(crawlerId, false)
				updatePendingTasksCounter(false)
			} else {
				crawlerLog(crawlerId, "Closing worker...")
				return
			}
		}
	}
}

// Crawl manager
func crawlManager(pendingTasks chan *tasks.Task, finishedTasks chan interface{}) {
	for {
		time.Sleep(100 * time.Millisecond)
		if newURLsPending() {
			urls := getPendingURLsPack(10)
			for _, url := range urls {
				pendingTasks <- tasks.NewTask(url)
				updatePendingTasksCounter(true)
			}
		} else if checkWorkersFinished() && taskCounterExhausted() {
			crawlerLogger.Printf("Closing all workers.")
			close(pendingTasks)
			time.Sleep(2000 * time.Millisecond)
			break
		}
	}
	finishedTasks <- 1
}

// Update worker status
func updateWorkerStatus(workerId int, workerStatus bool) {
	workerStatusMutex.Lock()
	defer workerStatusMutex.Unlock()
	workersStatus[workerId] = workerStatus
}

// Check if all workers have finished
func checkWorkersFinished() bool {
	workerStatusMutex.Lock()
	defer workerStatusMutex.Unlock()
	for _, workerStatus := range workersStatus {
		if workerStatus {
			return false
		}
	}
	return true
}

// Update pending tasks counter
func updatePendingTasksCounter(incr bool) {
	pendingTasksCounterMutex.Lock()
	defer pendingTasksCounterMutex.Unlock()
	if incr {
		pendingTasksCounter += 1
	} else {
		pendingTasksCounter -= 1
	}
}

// Check if tasks counter are finished
func taskCounterExhausted() bool {
	pendingTasksCounterMutex.Lock()
	defer pendingTasksCounterMutex.Unlock()
	return pendingTasksCounter == 0
}
