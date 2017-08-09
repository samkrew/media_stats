package db

import (
	"gopkg.in/mgo.v2"

	"time"
	"github.com/samkrew/media_stats/config"
	"sync"
	"github.com/samkrew/media_stats/logger"
)

var (
	mgoSession *mgo.Session
	mgoOnce    sync.Once
)

func connectMongoDb() *mgo.Session {
	dialInfo, err := mgo.ParseURL(*config.MongoURI)

	dialInfo.FailFast = false
	dialInfo.Timeout = 10 * time.Second

	session, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		logger.L.Fatal(err)
	}
	session.SetMode(mgo.Monotonic, true)
	return session
}

func GetInstance() (*mgo.Session, *mgo.Database) {
	mgoOnce.Do(func() {
		mgoSession = connectMongoDb()
	})
	s := mgoSession.Copy()

	return s, s.DB(*config.MongoDBName)
}
