package apigateway

import (
	"catch-all/gen/openapi"
	"strings"

	"github.com/labstack/echo/v4"
)

type PathParams = map[string]string

type Route struct {
	Handler echo.HandlerFunc
	Params  PathParams
}

type HandlerMap struct {
	Route   echo.Route
	Handler echo.HandlerFunc
}

func (m HandlerMap) match(path, method string) (bool, PathParams) {
	if m.Route.Method != method {
		return false, PathParams{}
	}

	if strings.Contains(m.Route.Path, ":") {
		params := PathParams{}
		strs1 := strings.Split(m.Route.Path, "/")
		strs2 := strings.Split(path, "/")
		if len(strs1) != len(strs2) {
			return false, PathParams{}
		}
		for i, w1 := range strs1 {
			w2 := strs2[i]
			if strings.Contains(w1, ":") {
				params[strings.TrimPrefix(w1, ":")] = w2
				continue
			}
			if w1 != w2 {
				return false, params
			}
		}

		return true, params
	} else {
		return m.Route.Path == path, PathParams{}
	}
}

type Router struct {
	Map []HandlerMap
}

func (r *Router) GetRoute(path, method string) *Route {
	for _, m := range r.Map {
		if b, params := m.match(path, method); b {
			return &Route{
				Handler: m.Handler,
				Params:  params,
			}
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
