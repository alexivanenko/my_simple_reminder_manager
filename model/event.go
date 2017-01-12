package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

type Event struct {
	ObjectID bson.ObjectId `json:"id" bson:"_id"`
	Name     string        `json:"name" bson:"name"`
	Date     time.Time     `json:"date" bson:"date"`
	TimeZone string        `json:"timezone" bson:"timezone"`
	ChatID   int64         `json:"chat_id" bson:"chat_id"`
}

func (event *Event) LoadAll() ([]Event, error) {
	db := GetDB()
	b := db.C("events")

	var result []Event
	err := b.Find(bson.M{}).All(&result)

	return result, err
}

func (event *Event) LoadCurrent() ([]Event, error) {
	db := GetDB()
	b := db.C("events")

	now := time.Now().UTC().Format("2006-01-02 15:04") + ":00"
	date, err := time.Parse("2006-01-02 15:04:05", now)

	var result []Event
	err = b.Find(bson.M{"date": date}).All(&result)

	return result, err
}

func (event *Event) Remove() error {
	db := GetDB()
	b := db.C("events")

	return b.Remove(bson.M{"_id": event.ObjectID})
}
