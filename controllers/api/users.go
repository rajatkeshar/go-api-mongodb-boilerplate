package api

import (
    "github.com/GORest-API-MongoDB/models"
    "github.com/GORest-API-MongoDB/controllers/handlers"
)

var UsersRoutes = models.RoutePrefix{
	"/api/users",
	[]models.Route{
		models.Route{
			"GetUsers",
			"GET",
			"",
			handlers.GetUsers,
			true,
		},
		models.Route{
			"GetUser",
			"GET",
			"/{userId}",
			handlers.GetUserById,
			true,
		},
		models.Route{
			"RegisterUser",
			"POST",
			"/register",
			handlers.RegisterNewUser,
			false,
		},
		models.Route{
			"DeleteUser",
			"DELETE",
			"/{userId}",
			handlers.DeleteUser,
			true,
		},
		models.Route{
			"UpdateUser",
			"PUT",
			"/{userId}",
			handlers.UpdateUser,
			true,
		},
	},
}
