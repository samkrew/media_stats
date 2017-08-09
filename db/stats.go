package db

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

type Stats struct {
	TaskId string `bson:"taskId" json:"taskId"`
	Url string `bson:"url" json:"url"`
	Hash string `bson:"hash" json:"hash"`
	Finished bool `bson:"finished" json:"finished"`
	Status string `bson:"status,omitempty" json:"status,omitempty"`
	Success bool `bson:"success" json:"success"`
	Resolution string `bson:"resolution,omitempty" json:"resolution,omitempty"`
	Bitrate string `bson:"bitrate,omitempty" json:"bitrate,omitempty"`
	Error string `bson:"error,omitempty" json:"error,omitempty"`
	Created time.Time `bson:"created" json:"created"`
}

func NewTaskStats(taskId string, url string, hash string) error {
	s,d := GetInstance()
	defer s.Close()

	return d.C("stats").Insert(Stats{
		TaskId: taskId,
		Url: url,
		Hash: hash,
		Finished: false,
		Created: time.Now(),
	})
}

func UpdateStatus(taskId string, status string) error {
	s,d := GetInstance()
	defer s.Close()

	return d.C("stats").Update(bson.M{"taskId": taskId}, bson.M{
		"$set": bson.M{
			"status": status,
		},
	})
}

func FinishStatsWithError(taskId string, error string) error {
	s,d := GetInstance()
	defer s.Close()

	return d.C("stats").Update(bson.M{"taskId": taskId}, bson.M{
		"$set": bson.M{
			"finished": true,
			"status": "File transfered",
			"success": false,
			"error": error,
		},
	})
}

func FinishStats(taskId string, resolution string, bitrate string) error {
	s,d := GetInstance()
	defer s.Close()

	return d.C("stats").Update(bson.M{"taskId": taskId}, bson.M{
		"$set": bson.M{
			"finished": true,
			"status": nil,
			"success": true,
			"resolution": resolution,
			"bitrate": bitrate,
		},
	})
}

func GetAllStats() (stats []*Stats) {
	s,d := GetInstance()
	defer s.Close()

	d.C("stats").Find(bson.M{}).Sort("-created").All(&stats)
	return
}

func init() {
	s, d := GetInstance()
	defer s.Close()

	d.C("stats").EnsureIndex(mgo.Index{
		Key:        []string{"taskId"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     false,
	})
}
