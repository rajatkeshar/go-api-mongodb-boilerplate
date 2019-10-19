package handlers

import (
    "os"
    //"log"
    "strings"
    "net/http"
    "encoding/json"
    "github.com/gorilla/mux"
    "golang.org/x/crypto/bcrypt"
    "github.com/go-api-mongodb-boilerplate/dao"
    "github.com/go-api-mongodb-boilerplate/models"
    "github.com/go-api-mongodb-boilerplate/lib/auth"
    "github.com/go-api-mongodb-boilerplate/lib/mailer"
    "github.com/go-api-mongodb-boilerplate/lib/responseHandler"
)

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

func UsersForgotPassword(w http.ResponseWriter, r *http.Request) {
    var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid request payload", false)
		return
	}

    user, err := dao.FindByEmailId(user.Email, user.Username)
    if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

    token, _ := auth.GenerateJWT(user)
    URL := "http://" + os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/api/auth/verify/forgot/password/" + token
    go mailer.SendMail(user.Email, "Forgot Password - go-api-mongodb-boilerplate", "Hi " + strings.Title(strings.ToLower(user.Firstname)) + "\nClick on the link below to reset your password \n" + "link: " + URL)
    response.Json(w, http.StatusCreated, "Reset Password Link Send On Your Registered Email", nil)
}

func UsersVerifyForgotPassword(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
    if err := auth.VerifyToken(w, params["token"]); err != nil {
        response.Json(w, http.StatusBadRequest, err.Error() + ", Reset Password Failed!", false)
        return
    }

    var user models.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid request payload", false)
		return
	}

    res, err := dao.FindById(w.Header().Get("_id"))
    if err != nil {
		response.Json(w, http.StatusBadRequest, "Invalid user ID", false)
		return
	}

    hash, _ := bcrypt.GenerateFromPassword([]byte(user.Password), 5)
	res.Password = string(hash)
    if err := dao.Update(res); err != nil {
		response.Json(w, http.StatusInternalServerError, err.Error(), false)
		return
	}
	response.Json(w, http.StatusOK, "Password Changed Success!", nil)
}
