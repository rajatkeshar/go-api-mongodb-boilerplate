package handlers

import (
        "os"
        //"log"
        "net/http"
        "strings"
        "encoding/json"
        "gopkg.in/mgo.v2/bson"
        "github.com/gorilla/mux"
        "golang.org/x/crypto/bcrypt"
        "github.com/go-api-mongodb-boilerplate/dao"
        "github.com/go-api-mongodb-boilerplate/models"
        "github.com/go-api-mongodb-boilerplate/lib/auth"
        "github.com/go-api-mongodb-boilerplate/lib/mailer"
        "github.com/go-api-mongodb-boilerplate/lib/multipart"
        "github.com/go-api-mongodb-boilerplate/lib/responseHandler"
)

// GET list of users
func GetUsers(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "Users Data", users)
}

// GET a users by its ID
func GetUserById(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["userId"])
	if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}
	response.Json(w, http.StatusOK, "User data", user)
}

// POST a new user
func RegisterNewUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid request payload", false)
		return
	}

	res, _ := dao.FindByEmailId(user.Email, user.Username)
	if res.Email != "" {
		response.Json(w, http.StatusBadRequest, "User Already Exists!", false)
		return
	}

	user.ID = bson.NewObjectId()
	hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	user.Password = string(hash)
	if err := dao.Insert(user); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
    token, _ := auth.GenerateJWT(user)
    URL := "http://" + os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/auth/verify/" + token
    user.Password = ""
    go mailer.SendMail(user.Email, "Verify Email - go-api-mongodb-boilerplate", "Hi " + strings.Title(strings.ToLower(user.Firstname)) + "\nRegistration Successful \n" + "Please Verify Your email: " + URL)
    response.Json(w, http.StatusCreated, "User Registration Success, Please Confirm Email", user)
}

// PUT update an existing user
func UpdateUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	params := mux.Vars(r)
	var user models.User

    user, err := dao.FindById(params["userId"])
	if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

	user.ID = bson.ObjectIdHex(params["userId"])
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid request payload", false)
		return
	}
	if err := dao.Update(user); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "User detail updated successfully", nil)
}

// DELETE an existing user
func DeleteUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	var user models.User
    params := mux.Vars(r)
    
    user, err := dao.FindById(params["userId"])
	if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}
	if err := dao.Delete(user); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "User deleted successfully", nil)
}

// UPLOAD profile for existing user
func UploadProfileForUser(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
    var user models.User
    params := mux.Vars(r)

    user, err := dao.FindById(params["userId"])
	if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

    data, err := multipart.UploadFile(w, r);
    if err != nil {
        response.Json(w, http.StatusInternalServerError, err.Error(), false)
    }

    user.ID = bson.ObjectIdHex(params["userId"])
    user.Url = data.(map[string]string)["filename"]

	if err := dao.Update(user); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "User profile updated successfully", nil)
}
