package dao

import (
	"log"
	"gopkg.in/mgo.v2"
	. "github.com/go-api-mongodb-boilerplate/config"
)

var config = Config{}
var db *mgo.Database

const (
	COLLECTION = "users"
)

// Parse the configuration file 'config.toml', and establish a connection to DB
func Connect() {
	config.Read()
	session, err := mgo.Dial(config.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(config.Database)
}

// Ensure Indexes on keys
func PopulateIndex() {
	for _, key := range []string{"username", "email"} {
		index := mgo.Index{
				Key:    []string{key},
				Unique: true,
		}
		db.C(COLLECTION).EnsureIndex(index);
	}
}
