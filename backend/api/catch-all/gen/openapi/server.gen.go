// Package openapi provides primitives to interact with the openapi HTTP API.
//
// Code generated by github.com/deepmap/oapi-codegen version v1.13.0 DO NOT EDIT.
package openapi

import (
	"fmt"
	"net/http"

	"github.com/deepmap/oapi-codegen/pkg/runtime"
	"github.com/gin-gonic/gin"
)

// ServerInterface represents all server handlers.
type ServerInterface interface {

	// (GET /storage/files/{key})
	GetFileUrl(c *gin.Context, key string)

	// (GET /storage/files/{key}/thumbnails)
	GetFileThumbnailUrls(c *gin.Context, key string)

	// (POST /transactions/storing_file)
	CreateTransaction(c *gin.Context)

	// (PATCH /transactions/storing_file/{transaction_id})
	UpdateTransaction(c *gin.Context, transactionId string)
}

// ServerInterfaceWrapper converts contexts to parameters.
type ServerInterfaceWrapper struct {
	Handler            ServerInterface
	HandlerMiddlewares []MiddlewareFunc
	ErrorHandler       func(*gin.Context, error, int)
}

type MiddlewareFunc func(c *gin.Context)

// GetFileUrl operation middleware
func (siw *ServerInterfaceWrapper) GetFileUrl(c *gin.Context) {

	var err error

	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameter("simple", false, "key", c.Param("key"), &key)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetFileUrl(c, key)
}

// GetFileThumbnailUrls operation middleware
func (siw *ServerInterfaceWrapper) GetFileThumbnailUrls(c *gin.Context) {

	var err error

	// ------------- Path parameter "key" -------------
	var key string

	err = runtime.BindStyledParameter("simple", false, "key", c.Param("key"), &key)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter key: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.GetFileThumbnailUrls(c, key)
}

// CreateTransaction operation middleware
func (siw *ServerInterfaceWrapper) CreateTransaction(c *gin.Context) {

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.CreateTransaction(c)
}

// UpdateTransaction operation middleware
func (siw *ServerInterfaceWrapper) UpdateTransaction(c *gin.Context) {

	var err error

	// ------------- Path parameter "transaction_id" -------------
	var transactionId string

	err = runtime.BindStyledParameter("simple", false, "transaction_id", c.Param("transaction_id"), &transactionId)
	if err != nil {
		siw.ErrorHandler(c, fmt.Errorf("Invalid format for parameter transaction_id: %s", err), http.StatusBadRequest)
		return
	}

	for _, middleware := range siw.HandlerMiddlewares {
		middleware(c)
		if c.IsAborted() {
			return
		}
	}

	siw.Handler.UpdateTransaction(c, transactionId)
}

// GinServerOptions provides options for the Gin server.
type GinServerOptions struct {
	BaseURL      string
	Middlewares  []MiddlewareFunc
	ErrorHandler func(*gin.Context, error, int)
}

// RegisterHandlers creates http.Handler with routing matching OpenAPI spec.
func RegisterHandlers(router gin.IRouter, si ServerInterface) {
	RegisterHandlersWithOptions(router, si, GinServerOptions{})
}

// RegisterHandlersWithOptions creates http.Handler with additional options
func RegisterHandlersWithOptions(router gin.IRouter, si ServerInterface, options GinServerOptions) {
	errorHandler := options.ErrorHandler
	if errorHandler == nil {
		errorHandler = func(c *gin.Context, err error, statusCode int) {
			c.JSON(statusCode, gin.H{"msg": err.Error()})
		}
	}

	wrapper := ServerInterfaceWrapper{
		Handler:            si,
		HandlerMiddlewares: options.Middlewares,
		ErrorHandler:       errorHandler,
	}

	router.GET(options.BaseURL+"/storage/files/:key", wrapper.GetFileUrl)
	router.GET(options.BaseURL+"/storage/files/:key/thumbnails", wrapper.GetFileThumbnailUrls)
	router.POST(options.BaseURL+"/transactions/storing_file", wrapper.CreateTransaction)
	router.PATCH(options.BaseURL+"/transactions/storing_file/:transaction_id", wrapper.UpdateTransaction)
}