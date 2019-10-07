# GoLang, Mux, MongoDB, JWT Restful API Boilerplate

This repository contains the web app to learn Go REST API development using Mux(Router), JSON Web Token (JWT), and MongoDB.

## Instructions

This is for educational purposes only and probably unsuitable for production

### Install Go Programming language latest version

[![N|Solid](https://sdtimes.com/wp-content/uploads/2018/02/golang.sh_-490x490.png)](https://golang.org/dl/)

### To get basic external modules for REST API

 ```sh
go get github.com/gorilla/mux gopkg.in/mgo.v2 github.com/BurntSushi/toml
```

* [mux](https://github.com/gorilla/mux) - Request router and dispatcher for matching incoming requests to their respective handler
* [mgo](https://gopkg.in/mgo.v2) - MongoDB driver
* [toml](https://github.com/BurntSushi/toml) - Parse the configuration file (MongoDB server & credentials)


### What's included

Basic CRUD routes for user management

Show Users GET      /api/users 
Show User GET       /api/users/{userId} 
Create User POST    /api/users/register 
Update User PUT     /api/users/{userId} 
Delete User DELETE  /api/users/{userId} 

User Login POST     /api/auth/login 
User Logout GET     /api/auth/logout 

Several routes are protected and require JWT tokens, which can be generated using the login route. You will need to create a user by sending a post request to the createUser route.

### Configuration

Database configuration will be done in config.toml file and rest of the configurating of project will be inside the .env file.

### To get this repository and run

 ```sh
$ git clone https://github.com/rajatkeshar/go-rest-api-jwt-mongo-boilerplate.git
$ go run *.go
```

## Todos

[] Implements Swagger on the top of middleware.
[] Implements Mail service while creating users.
[] Making the code more enhensive and moduler.
[] Implements redisDB for session management.
