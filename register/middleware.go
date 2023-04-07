package register

import "net/http"

type Middleware interface {
	Handle(next http.Handler) http.Handler
	GetName() string
	GetDescription() string
}

type MiddlewareRegister []interface{}
