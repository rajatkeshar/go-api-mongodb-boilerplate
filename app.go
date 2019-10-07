package main

import (
	"os"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/gorilla/handlers"
	. "github.com/GORest-API-MongoDB/dao"
	. "github.com/GORest-API-MongoDB/config"
)

var config = Config{}
var dao = UsersDAO{}

// Parse the configuration file 'config.toml', and establish a connection to DB
func init() {
	config.Read()
	dao.Server = config.Server
	dao.Database = config.Database
	dao.Connect()
	dao.PopulateIndex()
}

// Define HTTP request routes
func main() {
	routes := RoutesLoader()
	godotenv.Load()

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Server Is Running At ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), handlers.CORS(headersOk, originsOk, methodsOk)(routes)))
}
