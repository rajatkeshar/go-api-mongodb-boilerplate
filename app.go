package main

import (
	"os"
	"log"
	"net/http"
	"github.com/joho/godotenv"
	"github.com/gorilla/handlers"
	"github.com/go-api-mongodb-boilerplate/dao"
	"github.com/go-api-mongodb-boilerplate/router"
)

func init() {
	dao.Connect()
	dao.PopulateIndex()
}

// Define HTTP request routes
func main() {
	godotenv.Load()
	routes := router.NewRouter()

	originsOk := handlers.AllowedOrigins([]string{"*"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"})

	log.Println("Server Is Running At ", os.Getenv("PORT"))
	log.Fatal(http.ListenAndServe(":" + os.Getenv("PORT"), handlers.CORS(headersOk, originsOk, methodsOk)(routes)))
}
