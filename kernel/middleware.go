package kernel

import (
	"github.com/gorilla/mux"
	"github.com/krishpranav/gomicro/register"
)

func parseMiddleware(mwList []register.Middleware) []mux.MiddlewareFunc {
	var midFunc []mux.MiddlewareFunc

	for i := len(mwList) - 1; i > -1; i-- {
		midFunc = append(midFunc, mwList[i].Handle)
	}

	return midFunc
}
