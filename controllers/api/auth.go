package api

import (
    "github.com/GoRest-API-MongoDB-Boilerplate/models"
    "github.com/GoRest-API-MongoDB-Boilerplate/controllers/handlers"
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
	},
}
