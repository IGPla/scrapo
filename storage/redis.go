package storage

import (
	"fmt"
	"github.com/go-redis/redis"
	"github.com/IGPla/scrapo/logger"
	"github.com/IGPla/scrapo/tasks"
	"log"
	"os"
	"time"
)

var redisStorageLogger *log.Logger

func init() {
	redisStorageLogger = logger.GetLogger("REDISTORAGE", os.Stdout)
}

type RedisStorage struct {
	Host              string
	Port              int
	Password          string
	Db                int
	DefaultExpiration int
}

func (rs RedisStorage) StoreData(task *tasks.Task) error {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%v:%d", rs.Host, rs.Port),
		Password: rs.Password,
		DB:       rs.Db,
	})
	redisStorageLogger.Printf("Storing task in redis (%v)",
		task.URL)
	err := redisClient.Set(task.URL, task.Content, time.Duration(rs.DefaultExpiration)).Err()
	if err != nil {
		redisStorageLogger.Printf("An error arised storing data (%v)",
			err.Error())
		return err
	}
	return nil
}
