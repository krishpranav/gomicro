package kernel

import (
	"log"
	"net/http"
	"reflect"
)

type BaseController struct {
	Response http.ResponseWriter
	Request  *http.Request
}

var (
	BC BaseController
)

func registerBaseController(res http.ResponseWriter, req *http.Request, controller *interface{}) *interface{} {
	if err := setBaseController(res, req); err != nil {
		log.Fatal(err)
	}

	c := reflect.ValueOf(*controller).Elem().FieldByName("BaseController")
	c.Set(reflect.ValueOf(BC))

	return controller
}

func setBaseController(res http.ResponseWriter, req *http.Request) error {
	BC = BaseController{
		Response: res,
		Request:  req,
	}

	return nil
}
