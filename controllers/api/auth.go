package api

import (
    "github.com/go-api-mongodb-boilerplate/models"
    "github.com/go-api-mongodb-boilerplate/controllers/handlers"
)

var AuthRoutes = models.RoutePrefix{
	"/api/auth",
	[]models.Route{
		models.Route{
			"UsersLogin",
			"POST",
			"/login",
			handlers.UsersLogin,
			false,
		},
		models.Route{
			"UsersLogout",
			"GET",
			"/logout",
			handlers.UsersLogout,
			false,
		},
        models.Route{
			"UsersVerify",
			"GET",
			"/verify/{token}",
			handlers.UsersVerify,
			false,
		},
        models.Route{
			"UsersForgotPassword",
			"POST",
			"/forgot",
			handlers.UsersForgotPassword,
			false,
		},
        models.Route{
			"UsersVerifyForgotPassword",
			"POST",
			"/verify/forgot/password/{token}",
			handlers.UsersVerifyForgotPassword,
			false,
		},
	},
}
