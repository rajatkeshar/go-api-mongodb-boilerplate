package router

import (
    "os"
    "log"
    "regexp"
    "net/http"
    "github.com/gorilla/mux"
    "github.com/swaggo/http-swagger"
    _ "github.com/go-api-mongodb-boilerplate/docs"
    "github.com/go-api-mongodb-boilerplate/models"
    "github.com/go-api-mongodb-boilerplate/lib/auth"
    "github.com/go-api-mongodb-boilerplate/controllers/api"
)

func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        log.Println(r.RequestURI)
        // Call the next handler, which can be another middleware in the chain, or the final handler.
        next.ServeHTTP(w, r)
    })
}

func commonMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if match, _ := regexp.MatchString("/swagger/*", r.URL.Path); match {
            next.ServeHTTP(w, r)
            return
        }
        w.Header().Set("Content-Type", "application/json")
        next.ServeHTTP(w, r)
    })
}

func NewRouter() *mux.Router {
    routes := mux.NewRouter()

    routes.Use(loggingMiddleware)
	routes.Use(commonMiddleware)

    routes.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://" + os.Getenv("HOST") + ":" + os.Getenv("PORT") + "/swagger/doc.json"), //The url pointing to API definition"
	))

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
