package dao

import (
	"gopkg.in/mgo.v2/bson"
	"github.com/go-api-mongodb-boilerplate/models"
)

// Find list of users
func FindAll() ([]models.User, error) {
	var users []models.User
	err := db.C(COLLECTION).Find(bson.M{}).All(&users)
	return users, err
}

// Find a user by its id
func FindById(id string) (models.User, error) {
	var user models.User
	err := db.C(COLLECTION).FindId(bson.ObjectIdHex(id)).One(&user)
	return user, err
}

// Insert a user into database
func Insert(user models.User) error {
	err := db.C(COLLECTION).Insert(&user)
	return err
}

// Delete an existing user
func Delete(user models.User) error {
	err := db.C(COLLECTION).Remove(&user)
	return err
}

// Update an existing user
func Update(user models.User) error {
	err := db.C(COLLECTION).UpdateId(user.ID, &user)
	return err
}

// Find a user by its email id
func FindByEmailId(email string, username string) (models.User, error) {
	var user models.User
	err := db.C(COLLECTION).Find(bson.M{"$or": []interface{}{bson.M{"email": email}, bson.M{"username": username}}}).One(&user)
	return user, err
}
