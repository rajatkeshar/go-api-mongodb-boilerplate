package dao

import (
	"log"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/go-api-mongodb-boilerplate/models"
)

type UsersDAO struct {
	Server   string
	Database string
}

var db *mgo.Database

const (
	COLLECTION = "users"
)

// Establish a connection to database
func (m *UsersDAO) Connect() {
	session, err := mgo.Dial(m.Server)
	if err != nil {
		log.Fatal(err)
	}
	db = session.DB(m.Database)
}

// Ensure Indexes on keys
func (m *UsersDAO) PopulateIndex() {
	for _, key := range []string{"username", "email"} {
		index := mgo.Index{
				Key:    []string{key},
				Unique: true,
		}
		db.C(COLLECTION).EnsureIndex(index);
	}
}

// Find list of users
func (m *UsersDAO) FindAll() ([]models.User, error) {
	var users []models.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

// Find a user by its id
func (m *UsersDAO) FindById(id string) (models.User, error) {
	var user models.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Insert a user into database
func (m *UsersDAO) Insert(user models.User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete an existing user
func (m *UsersDAO) Delete(user models.User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update an existing user
func (m *UsersDAO) Update(user models.User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}

// Find a user by its email id
func (m *UsersDAO) FindByEmailId(email string, username string) (models.User, error) {
	var user models.User
	err := db.C(COLLECTION).Find(bson.M{"$or": []interface{}{bson.M{"email": email}, bson.M{"username": username}}}).One(&user)
	return user, err
}
