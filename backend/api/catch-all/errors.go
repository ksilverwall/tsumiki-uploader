package main

import "catch-all/gen/openapi"

type ErrorResponse struct {
	Status int
	Body   openapi.Error
}

func NewErrorResponse(err error) ErrorResponse {
	return ErrorResponse{
		Status: 500,
		Body: openapi.Error{
			Code:    500,
			Message: err.Error(),
		},
	}
}
