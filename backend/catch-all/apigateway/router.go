package apigateway

import (
	"catch-all/gen/openapi"

	"github.com/labstack/echo/v4"
)

type Route struct {
	Handler echo.HandlerFunc
}

type HandlerMap struct {
	Route   echo.Route
	Handler echo.HandlerFunc
}

type Router struct {
	Map []HandlerMap
}

func (r *Router) GetRoute(path, method string) *Route {
	for _, m := range r.Map {
		if m.Route.Path == path && m.Route.Method == method {
			return &Route{Handler: m.Handler}
		}
	}

	return nil
}

func (r *Router) Register(s openapi.ServerInterface) {
	openapi.RegisterHandlers(RouteRegister{Router: r}, s)
}

func NewRouter() Router {
	return Router{[]HandlerMap{}}
}
