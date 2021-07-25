package nonapi

import (
	"net/http"
)

type ExampleContext struct {
	Err *AppError
}

type AppError struct{}

func (er *AppError) Error() string {
	return "This is an error"
}

func iR() *AppError {
	return &AppError{}
}

func a(c *ExampleContext, w http.ResponseWriter, r *http.Request) {
	err := iR()
	if err != nil {
		c.Err = err
	}
}
