# GoLang Restful API Boilerplate (Mux, JWT, MongoDB)  

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

Basic CRUD routes for user management <br/>

Show Users GET      /api/users <br />
Show User GET       /api/users/{userId} <br /> 
Create User POST    /api/users/register <br />
Update User PUT     /api/users/{userId} <br />
Delete User DELETE  /api/users/{userId} <br />

User Login POST     /api/auth/login <br />
User Logout GET     /api/auth/logout <br />

Several routes are protected and require JWT tokens, which can be generated using the login route. You will need to create a user by sending a post request to the createUser route.

### Configuration

Database configuration will be done in config.toml file and rest of the configurating of project will be inside the .env file.

### To get this repository and run

 ```sh
$ git clone https://github.com/rajatkeshar/GoRest-API-MongoDB-Boilerplate.git
$ go run *.go
```

## Todos

[] Implements Swagger for API crud. <br />
[] Implements Mail service while creating users. <br />
[] Making the code more enhensive and moduler. <br />
[] Implements redisDB for session management. <br />
