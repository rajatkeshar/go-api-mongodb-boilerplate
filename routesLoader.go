package main

import (
    "log"
    "net/http"
    "github.com/gorilla/mux"
	//"github.com/bitly/go-simplejson"
    "github.com/GORest-API-MongoDB/models"
    "github.com/GORest-API-MongoDB/lib/auth"
    "github.com/GORest-API-MongoDB/controllers/api"
)

//Home Page
// func homePage(w http.ResponseWriter, r *http.Request) {
// 		jsonBuilder := simplejson.New()
// 		jsonBuilder.Set("_id", w.Header().Get("_id"))
// 		jsonBuilder.Set("firstname", w.Header().Get("firstname"))
// 		jsonBuilder.Set("lastname", w.Header().Get("lastname"))
// 		respondWithJson(w, http.StatusOK, "Home Page!", jsonBuilder)
// }

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

func RoutesLoader() *mux.Router {
    routes := mux.NewRouter()

    routes.Use(loggingMiddleware)
	routes.Use(commonMiddleware)
	// routes.Handle("/", auth.IsAuthorized(homePage)).Methods("GET")

    //append applications routes
    models.Routes = append(models.Routes, api.AuthRoutes)
	models.Routes = append(models.Routes, api.UsersRoutes)

	for _, route := range models.Routes {
		//create subroute
		routePrefix := routes.PathPrefix(route.Prefix).Subrouter()
		//loop through each sub route
		for _, r := range route.SubRoutes {
			var handler http.Handler
			handler = r.HandlerFunc
			//check to see if route should be protected with jwt
			if r.Protected {
				handler = auth.IsAuthorized(r.HandlerFunc)
			}
			//attach sub route
			routePrefix.
				Path(r.Pattern).
				Handler(handler).
				Methods(r.Method).
				Name(r.Name)
		}
	}
    return routes
}