package handlers

import (
        //"log"
        "net/http"
        "encoding/json"
        "gopkg.in/mgo.v2/bson"
        "github.com/gorilla/mux"
        "golang.org/x/crypto/bcrypt"
        . "github.com/GORest-API-MongoDB/dao"
        "github.com/GORest-API-MongoDB/models"
)

var dao = UsersDAO{}

// GET list of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "Users Data", users)
}

// GET a users by its ID
func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJson(w, http.StatusOK, "User data", user)
}

// POST a new user
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	res, _ := dao.FindByEmailId(user.Email, user.Username)
	if res.Email != "" {
		respondWithError(w, http.StatusBadRequest, "User Already Exists!")
		return
	}

	user.ID = bson.NewObjectId()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hash)
	if err := dao.Insert(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusCreated, "User created successfully", user)
}

// PUT update an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var user models.User
	user.ID = bson.ObjectIdHex(params["id"])
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Update(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "User detail updated successfully", "")
}

// DELETE an existing user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	if err := dao.Delete(user); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "User deleted successfully", "")
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bson.M{"code": code, "success": false, "msg": msg, "data": nil})
}

func respondWithJson(w http.ResponseWriter, code int, msg string, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bson.M{"code": code, "success": true, "msg": msg, "data": payload})
}
