package main

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"encoding/json"
	"gopkg.in/mgo.v2/bson"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
	"github.com/gorilla/handlers"
	"github.com/bitly/go-simplejson"
	. "github.com/GORest-API-MongoDB/dao"
	"github.com/GORest-API-MongoDB/models"
	. "github.com/GORest-API-MongoDB/config"
	auth "github.com/GORest-API-MongoDB/lib"
)

var config = Config{}
var dao = UsersDAO{}

//Home Page
func homePage(w http.ResponseWriter, r *http.Request) {
		jsonBuilder := simplejson.New()
		jsonBuilder.Set("_id", w.Header().Get("_id"))
		jsonBuilder.Set("firstname", w.Header().Get("firstname"))
		jsonBuilder.Set("lastname", w.Header().Get("lastname"))
		respondWithJson(w, http.StatusOK, "Home Page!", jsonBuilder)
}

// GET list of users
func AllUsersEndPoint(w http.ResponseWriter, r *http.Request) {
	users, err := dao.FindAll()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJson(w, http.StatusOK, "Users Data", users)
}

// GET a users by its ID
func FindUserEndpoint(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	user, err := dao.FindById(params["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid user ID")
		return
	}
	respondWithJson(w, http.StatusOK, "User data", user)
}

// POST a new user
func RegisterUserEndPoint(w http.ResponseWriter, r *http.Request) {
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
func UpdateUserEndPoint(w http.ResponseWriter, r *http.Request) {
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
func DeleteUserEndPoint(w http.ResponseWriter, r *http.Request) {
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

func LoginEndPoint(w http.ResponseWriter, r *http.Request) {
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

func LogoutEndPoint(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("_id", "")
	w.Header().Set("email", "")
	w.Header().Set("firstname", "")
	w.Header().Set("lastname", "")
	respondWithJson(w, http.StatusOK, "Logout Success!", nil)
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bson.M{"code": code, "success": false, "msg": msg, "data": nil})
}

func respondWithJson(w http.ResponseWriter, code int, msg string, payload interface{}) {
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(bson.M{"code": code, "success": true, "msg": msg, "data": payload})
}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()

	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
	dao.PopulateIndex()
}

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Do stuff here
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

// Define HTTP request routes
func main() {
	routes := mux.NewRouter()
	godotenv.Load()

	routes.Use(loggingMiddleware)
	routes.Use(commonMiddleware)

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	routes.Handle("/api/users", auth.IsAuthorized(AllUsersEndPoint)).Methods("GET")
	routes.Handle("/api/users/register", auth.IsAuthorized(RegisterUserEndPoint)).Methods("POST")
	routes.Handle("/api/users/{id}", auth.IsAuthorized(UpdateUserEndPoint)).Methods("PUT")
	routes.Handle("/api/users", auth.IsAuthorized(DeleteUserEndPoint)).Methods("DELETE")
	routes.Handle("/api/users/{id}", auth.IsAuthorized(FindUserEndpoint)).Methods("GET")
	routes.Handle("/api/auth/login", auth.IsAuthorized(LoginEndPoint)).Methods("POST")
	routes.Handle("/api/auth/logout", auth.IsAuthorized(LogoutEndPoint)).Methods("GET")
	routes.Handle("/", auth.IsAuthorized(homePage)).Methods("GET")

	fmt.Println("Server Is Running At ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), handlers.CORS(headersOk, originsOk, methodsOk)(routes)))
}
