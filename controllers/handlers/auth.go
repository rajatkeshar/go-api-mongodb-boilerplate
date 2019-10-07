package handlers

import (
    //"log"
    "net/http"
    "encoding/json"
    "golang.org/x/crypto/bcrypt"
    //. "github.com/GORest-API-MongoDB/dao"
    "github.com/GORest-API-MongoDB/models"
    "github.com/GORest-API-MongoDB/lib/auth"
    "github.com/GORest-API-MongoDB/lib/responseHandler"
)

//var dao = UsersDAO{}

func UsersLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid request payload", false)
		return
	}
	password := user.Password
	user, err := dao.FindByEmailId(user.Email, user.Username)
	if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid password!", false)
		return
	}

	token, _ := auth.GenerateJWT(user)
	w.Header().Set("Token", token)
	response.Json(w, http.StatusOK, "Login Success!", user)
}

func UsersLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("_id", "")
	w.Header().Set("email", "")
	w.Header().Set("firstname", "")
	w.Header().Set("lastname", "")
	response.Json(w, http.StatusOK, "Logout Success!", nil)
}
