package handlers

import (
    //"log"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    //. "github.com/GoRest-API-MongoDB-Boilerplate/dao"
    "github.com/GoRest-API-MongoDB-Boilerplate/models"
    "github.com/GoRest-API-MongoDB-Boilerplate/lib/auth"
    "github.com/GoRest-API-MongoDB-Boilerplate/lib/responseHandler"
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

    if !user.Verify {
        response.Json(w, http.StatusBadRequest, "User not verified", false)
		return
    }

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid password!", false)
		return
	}
    user.Password = ""
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

func UsersVerify(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    if err := auth.VerifyToken(w, params["token"]); err != nil {
        response.Json(w, http.StatusBadRequest, err.Error() + ", Email Verification Failed!", false)
        return
    }
    user, err := dao.FindById(w.Header().Get("_id"))

    if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

    if user.Verify {
        response.Json(w, http.StatusBadRequest, "User Already Verified!", false)
		return
    }
    user.Verify = true
    if err := dao.Update(user); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "Email Verified Success!", nil)
}
