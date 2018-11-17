package storage

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/IGPla/scrapo/tasks"
	"os/exec"
	"strings"
	"testing"
)

var redisHost string = "0.0.0.0"
var redisPort int = 12345
var redisDb int = 0
var redisDefaultExpiration = 0
var redisContainerName = "redistest"

func getDockerBinary() string {
	path, err := exec.LookPath("docker")
	if err != nil {
		panic(err)
	}
	return path
}

func openRedis() {
	cmdString := strings.Split(fmt.Sprintf("%v run --name %v -p %v:%d:6379 -d redis",
		getDockerBinary(), redisContainerName, redisHost, redisPort), " ")

	cmd := exec.Command(cmdString[0], cmdString[1:]...)
	err := cmd.Run()
	if err != nil {
		panic(err)
	}
}

func closeRedis() {
	cmds := []string{fmt.Sprintf("%v stop %v", getDockerBinary(),
		redisContainerName), fmt.Sprintf("%v rm %v", getDockerBinary(),
		redisContainerName)}
	for _, cmdStr := range cmds {
		cmdString := strings.Split(cmdStr, " ")
		cmd := exec.Command(cmdString[0], cmdString[1:]...)
		err := cmd.Run()
		if err != nil {
			panic(err)
		}
	}
}

func getRedisConnection() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%d", redisHost, redisPort),
		Password: "",
		DB:       redisDb,
	})
}

func TestRedisStoreData(t *testing.T) {
	// Base content
	var content string = "Test content"
	var contentBytes = []byte(content)
	var task *tasks.Task = new(tasks.Task)
	task.URL = "http://www.google.com/test1/"
	task.Content = contentBytes
	task.ContentType = "text/html; charset=UTF-8"

	// Redis launch (docker)
	openRedis()
	defer closeRedis()
	redisClient := getRedisConnection()

	// Redisstorage init
	var rs *RedisStorage = new(RedisStorage)
	rs.Host = "localhost"
	rs.Port = redisPort

	// Store file without name
	htmlError := rs.StoreData(task)
	if htmlError != nil {
		t.Errorf("Errors appeared while storing content. (%v)",
			htmlError.Error())
	}
	val1, getErr1 := redisClient.Get(task.URL).Result()
	if getErr1 != nil {
		t.Errorf("Error arised while retrieving stored object (%v) from redis. (%v)",
			task.URL, getErr1.Error())
	}
	if val1 != content {
		t.Errorf("Wrong retrieved value. e(%v), g(%v)",
			content,
			val1)
	}
	// Store file with name
	task.URL = "http://www.google.com/test2?p1=v1"
	htmlError = rs.StoreData(task)
	if htmlError != nil {
		t.Errorf("Errors appeared while storing content. (%v)",
			htmlError.Error())
	}
	val2, getErr2 := redisClient.Get(task.URL).Result()
	if getErr2 != nil {
		t.Errorf("Error arised while retrieving stored object (%v) from redis. (%v)",
			task.URL, getErr2.Error())
	}
	if val2 != content {
		t.Errorf("Wrong retrieved value. e(%v), g(%v)",
			content,
			val2)
	}

}
