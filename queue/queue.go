package queue

import (
	"fmt"
	"github.com/tsuru/monsterqueue"
	"github.com/tsuru/monsterqueue/mongodb"
	"sync"
	"time"
	"github.com/samkrew/media_stats/config"
)

type Queue struct {
	sync.Mutex
	Instance monsterqueue.Queue
}

var (
	clientQueue Queue
	serverQueue Queue
)

func (q *Queue) Reset() {
	q.Lock()
	defer q.Unlock()
	if q.Instance != nil {
		q.Instance.Stop()
		q.Instance.ResetStorage()
		q.Instance = nil
	}
}

func (q *Queue) Stop() {
	q.Lock()
	defer q.Unlock()
	if q.Instance != nil {
		q.Instance.Stop()
	}
}

func (q *Queue) Run() {
	q.Lock()
	defer q.Unlock()
	if q.Instance != nil {
		q.Instance.ProcessLoop()
	}
}

func GetServerQueue() (*Queue, error) {
	serverQueue.Lock()
	defer serverQueue.Unlock()
	if serverQueue.Instance != nil {
		return &serverQueue, nil
	}
	queueMongoUrl := *config.MongoURI
	queueMongoDB := *config.MongoDBName

	conf := mongodb.QueueConfig{
		CollectionPrefix: "grabber",
		Url:              queueMongoUrl,
		Database:         queueMongoDB,
		PollingInterval:  time.Duration(time.Second),
	}
	var err error
	serverQueue.Instance, err = mongodb.NewQueue(conf)
	if err != nil {
		return nil, fmt.Errorf("Could not create serverQueue instance. Error: %s", err)
	}

	return &serverQueue, nil
}

func GetClientQueue() (*Queue, error) {
	clientQueue.Lock()
	defer clientQueue.Unlock()
	if clientQueue.Instance != nil {
		return &clientQueue, nil
	}
	queueMongoUrl := *config.MongoURI
	queueMongoDB := *config.MongoDBName

	conf := mongodb.QueueConfig{
		CollectionPrefix: "grabber",
		Url:              queueMongoUrl,
		Database:         queueMongoDB,
		PollingInterval:  time.Duration(time.Second),
	}
	var err error
	clientQueue.Instance, err = mongodb.NewQueue(conf)
	if err != nil {
		return nil, fmt.Errorf("Could not create serverQueue instance. Error: %s", err)
	}

	return &clientQueue, nil
}