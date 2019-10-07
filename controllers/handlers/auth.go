package handlers

import (
    //"log"
    "net/http"
    "encoding/json"
    //"gopkg.in/mgo.v2/bson"
    "golang.org/x/crypto/bcrypt"
    //. "github.com/GORest-API-MongoDB/dao"
    "github.com/GORest-API-MongoDB/models"
    "github.com/GORest-API-MongoDB/lib/auth"
)

//var dao = UsersDAO{}

func UsersLogin(w http.ResponseWriter, r *http.Request) {
	var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	password := user.Password
	user, err := dao.FindByEmailId(user.Email, user.Username)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid password!")
		return
	}

	token, _ := auth.GenerateJWT(user)
	w.Header().Set("Token", token)
	respondWithJson(w, http.StatusOK, "Login Success!", user)
}

func UsersLogout(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("_id", "")
	w.Header().Set("email", "")
	w.Header().Set("firstname", "")
	w.Header().Set("lastname", "")
	respondWithJson(w, http.StatusOK, "Logout Success!", nil)
}

// func respondWithError(w http.ResponseWriter, code int, msg string) {
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(bson.M{"code": code, "success": false, "msg": msg, "data": nil})
// }
//
// func respondWithJson(w http.ResponseWriter, code int, msg string, payload interface{}) {
// 	w.WriteHeader(code)
// 	json.NewEncoder(w).Encode(bson.M{"code": code, "success": true, "msg": msg, "data": payload})
// }
