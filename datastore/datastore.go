package datastore

import (
	"github.com/globalsign/mgo"
	"github.com/sirupsen/logrus"
)

type DataStore struct {
	Session *mgo.Session
}

func CreateStore() *DataStore {
	session, err := mgo.Dial("mongodb://mongo:27017")

	if err != nil {
		logrus.Panicf("Error connecting to mongo: %s", err)
	}

	return &DataStore{
		Session: session,
	}
}