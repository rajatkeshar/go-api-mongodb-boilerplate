package models

import "net/http"

var Routes []RoutePrefix

type RoutePrefix struct {
	Prefix string
	SubRoutes []Route
}

type Route struct {
	Name string
	Method string
	Pattern string
	HandlerFunc http.HandlerFunc
	Protected bool
}
