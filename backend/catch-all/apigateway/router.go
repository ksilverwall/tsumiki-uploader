package apigateway

import (
	"catch-all/gen/openapi"
	"strings"

	"github.com/labstack/echo/v4"
)

type Route struct {
	Handler echo.HandlerFunc
}

type HandlerMap struct {
	Route   echo.Route
	Handler echo.HandlerFunc
}

func (m HandlerMap) match(path, method string) bool {
	if m.Route.Method != method {
		return false
	}

	if strings.Contains(m.Route.Path, ":") {
		strs1 := strings.Split(m.Route.Path, "/")
		strs2 := strings.Split(path, "/")
		if len(strs1) != len(strs2) {
			return false
		}
		for i, w1 := range strs1 {
			if strings.Contains(w1, ":") {
				continue
			}
			w2 := strs2[i]
			if w1 != w2 {
				return false
			}
		}

		return true
	} else {
		return m.Route.Path == path
	}
}

type Router struct {
	Map []HandlerMap
}

func (r *Router) GetRoute(path, method string) *Route {
	for _, m := range r.Map {
		if m.match(path, method) {
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
