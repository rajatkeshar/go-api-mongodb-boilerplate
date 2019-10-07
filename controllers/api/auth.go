package api

import (
    "github.com/GORest-API-MongoDB/models"
    "github.com/GORest-API-MongoDB/controllers/handlers"
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
	},
}
