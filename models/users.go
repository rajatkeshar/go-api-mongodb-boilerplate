package models

import "gopkg.in/mgo.v2/bson"

// Represents a user, we uses bson keyword to tell the mgo driver how to name
// the properties in mongodb document
type User struct {
	ID bson.ObjectId `bson:"_id" json:"_id"`
	Firstname string `bson:"firstname" json:"firstname"`
	Lastname string `bson:"lastname" json:"lastname"`
	Age  int `bson:"age" json:"age"`
	Email string `bson:"email" json:"email"`
	Username string `bson:"username" json:"username"`
	Password string `bson:"password,omitempty" json:"password,omitempty"`
	Verify bool `bson:"verify" json:"verify" default:"false"`
	Url string `bson:"url,omitempty" json:"url,omitempty"`
}
