package apigateway

import "github.com/labstack/echo/v4"

type RouteRegister struct {
	Router *Router
}

//
// implements EchoRouter methods
//

func (r RouteRegister) CONNECT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) DELETE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) GET(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	rt := echo.Route{
		Method: echo.GET,
		Path:   path,
		Name:   "",
	}

	r.Router.Map = append(r.Router.Map, HandlerMap{
		Route:   rt,
		Handler: h,
	})

	return &rt
}
func (r RouteRegister) HEAD(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) OPTIONS(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) PATCH(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) POST(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	rt := echo.Route{
		Method: echo.POST,
		Path:   path,
		Name:   "",
	}

	r.Router.Map = append(r.Router.Map, HandlerMap{
		Route:   rt,
		Handler: h,
	})

	return &rt
}
func (r RouteRegister) PUT(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
func (r RouteRegister) TRACE(path string, h echo.HandlerFunc, m ...echo.MiddlewareFunc) *echo.Route {
	return nil
}
