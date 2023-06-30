package apigateway

import (
	"encoding/json"
	"errors"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"

	"github.com/aws/aws-lambda-go/events"
	"github.com/labstack/echo/v4"
)

type ContextStatus struct {
	Response *events.APIGatewayProxyResponse
}

type Context struct {
	Status *ContextStatus
}

//
// Implement echo.Context
//

func (c Context) Request() *http.Request          { c.non_implemented_method_call(); return nil }
func (c Context) SetRequest(r *http.Request)      { c.non_implemented_method_call() }
func (c Context) SetResponse(r *echo.Response)    { c.non_implemented_method_call() }
func (c Context) Response() *echo.Response        { c.non_implemented_method_call(); return nil }
func (c Context) IsTLS() bool                     { c.non_implemented_method_call(); return false }
func (c Context) IsWebSocket() bool               { c.non_implemented_method_call(); return false }
func (c Context) Scheme() string                  { c.non_implemented_method_call(); return "" }
func (c Context) RealIP() string                  { c.non_implemented_method_call(); return "" }
func (c Context) Path() string                    { c.non_implemented_method_call(); return "" }
func (c Context) SetPath(p string)                { c.non_implemented_method_call(); return }
func (c Context) Param(name string) string        { c.non_implemented_method_call(); return "" }
func (c Context) ParamNames() []string            { c.non_implemented_method_call(); return make([]string, 0) }
func (c Context) SetParamNames(names ...string)   { c.non_implemented_method_call(); return }
func (c Context) ParamValues() []string           { c.non_implemented_method_call(); return make([]string, 0) }
func (c Context) SetParamValues(values ...string) { c.non_implemented_method_call() }
func (c Context) QueryParam(name string) string   { c.non_implemented_method_call(); return "" }
func (c Context) QueryParams() url.Values         { c.non_implemented_method_call(); return url.Values{} }
func (c Context) QueryString() string             { c.non_implemented_method_call(); return "" }
func (c Context) FormValue(name string) string    { c.non_implemented_method_call(); return "" }
func (c Context) FormParams() (url.Values, error) {
	c.non_implemented_method_call()
	return url.Values{}, errors.New("not implemented")
}
func (c Context) FormFile(name string) (*multipart.FileHeader, error) {
	c.non_implemented_method_call()
	return nil, errors.New("not implemented")
}
func (c Context) MultipartForm() (*multipart.Form, error) {
	c.non_implemented_method_call()
	return nil, errors.New("not implemented")
}
func (c Context) Cookie(name string) (*http.Cookie, error) {
	c.non_implemented_method_call()
	return nil, errors.New("not implemented")
}
func (c Context) SetCookie(cookie *http.Cookie) { c.non_implemented_method_call() }
func (c Context) Cookies() []*http.Cookie {
	c.non_implemented_method_call()
	return make([]*http.Cookie, 0)
}
func (c Context) Get(key string) interface{}      { c.non_implemented_method_call(); return nil }
func (c Context) Set(key string, val interface{}) { c.non_implemented_method_call() }
func (c Context) Bind(i interface{}) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Validate(i interface{}) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Render(code int, name string, data interface{}) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) HTML(code int, html string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) HTMLBlob(code int, b []byte) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) String(code int, s string) error {
	c.Status.Response = &events.APIGatewayProxyResponse{
		StatusCode: code,
		Body:       s,
	}

	return nil
}
func (c Context) JSON(code int, i interface{}) error {
	b, err := json.Marshal(i)
	if err != nil {
		return err
	}

	c.String(code, string(b))

	return nil
}
func (c Context) JSONPretty(code int, i interface{}, indent string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) JSONBlob(code int, b []byte) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) JSONP(code int, callback string, i interface{}) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) JSONPBlob(code int, callback string, b []byte) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) XML(code int, i interface{}) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) XMLPretty(code int, i interface{}, indent string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) XMLBlob(code int, b []byte) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Blob(code int, contentType string, b []byte) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Stream(code int, contentType string, r io.Reader) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) File(file string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Attachment(file string, name string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Inline(file string, name string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) NoContent(code int) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Redirect(code int, url string) error {
	c.non_implemented_method_call()
	return errors.New("not implemented")
}
func (c Context) Error(err error) { c.non_implemented_method_call() }
func (c Context) Handler() echo.HandlerFunc {
	c.non_implemented_method_call()
	panic("not implemented")
}
func (c Context) SetHandler(h echo.HandlerFunc)                { c.non_implemented_method_call() }
func (c Context) Logger() echo.Logger                          { c.non_implemented_method_call(); panic("not implemented") }
func (c Context) SetLogger(l echo.Logger)                      { c.non_implemented_method_call() }
func (c Context) Echo() *echo.Echo                             { c.non_implemented_method_call(); return nil }
func (c Context) Reset(r *http.Request, w http.ResponseWriter) { c.non_implemented_method_call() }

// Private Methods
func (c Context) non_implemented_method_call() {}

//
// Constructor
//

func NewContext() Context {
	return Context{
		Status: &ContextStatus{
			Response: &events.APIGatewayProxyResponse{StatusCode: 200, Body: "Default"},
		},
	}
}
