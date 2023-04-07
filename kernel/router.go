package kernel

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"reflect"
	"strings"

	"github.com/gorilla/mux"
	"github.com/krishpranav/gomicro/register"
	"github.com/krishpranav/gomicro/tool"
)

type Request map[string]interface{}

func WebRouter(routes []register.HTTPRouter) *mux.Router {
	router := mux.NewRouter()
	router.Use(gzipMiddleware)

	for _, r := range routes {
		if len(r.Route) > 0 {
			HandleSingleRoute(r.Route, router)
		}

		if len(r.Groups) > 0 {
			HandleGroups(r.Groups, router)
		}

		GiveAccessToPublicFolder(router)
	}

	return router
}

func HandleSingleRoute(routes []register.Route, router *mux.Router) {
	for _, route := range routes {
		hasMiddleware := len(route.Middleware) > 0
		directive := strings.Split(route.Action, "@")
		validation := route.Validation
		if hasMiddleware {
			subRouter := mux.NewRouter()
			subRouter.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
				if err := validateRequest(validation, request); err != nil {
					writer.WriteHeader(http.StatusUnprocessableEntity)
					_, _ = writer.Write([]byte(err.Error()))

					return
				}

				executeControllerDirective(directive, writer, request, validation)
			}).Methods(route.Method)

			subRouter.Use(parseMiddleware(route.Middleware)...)
			router.Handle(route.Path, subRouter).Methods(route.Method)
		} else {
			router.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
				if err := validateRequest(validation, request); err != nil {
					writer.WriteHeader(http.StatusUnprocessableEntity)
					_, _ = writer.Write([]byte(err.Error()))

					return
				}

				executeControllerDirective(directive, writer, request, validation)
			}).Methods(route.Method)
		}
	}
}

func HandleGroups(groups []register.Group, router *mux.Router) {
	for _, group := range groups {
		subRouter := router.PathPrefix(group.Prefix).Subrouter()

		for _, route := range group.Routes {
			directive := strings.Split(route.Action, "@")
			validation := route.Validation
			if len(route.Middleware) > 0 {
				nestedRouter := mux.NewRouter()
				fullPath := fmt.Sprintf("%s%s", group.Prefix, route.Path)
				nestedRouter.HandleFunc(fullPath, func(writer http.ResponseWriter, request *http.Request) {
					if err := validateRequest(validation, request); err != nil {
						writer.WriteHeader(http.StatusUnprocessableEntity)
						_, _ = writer.Write([]byte(err.Error()))

						return
					}

					executeControllerDirective(directive, writer, request, validation)
				}).Methods(route.Method)

				nestedRouter.Use(parseMiddleware(route.Middleware)...)
				subRouter.Handle(route.Path, nestedRouter).Methods(route.Method)
			} else {
				subRouter.HandleFunc(route.Path, func(writer http.ResponseWriter, request *http.Request) {
					if err := validateRequest(validation, request); err != nil {
						writer.WriteHeader(http.StatusUnprocessableEntity)
						_, _ = writer.Write([]byte(err.Error()))

						return
					}

					executeControllerDirective(directive, writer, request, validation)
				}).Methods(route.Method)
			}
		}

		subRouter.Use(parseMiddleware(group.Middleware)...)
	}
}

func GiveAccessToPublicFolder(router *mux.Router) {
	publicDirectory := http.Dir(tool.GetDynamicPath("public"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", http.FileServer(publicDirectory)))
}

func getControllerItem(itemName string) register.ControllerRegisterItem {
	var result register.ControllerRegisterItem
	for _, contr := range Controllers {
		controllerName := reflect.Indirect(reflect.ValueOf(contr.Controller)).Type().Name()
		if controllerName == itemName {
			result = contr
		}
	}

	return result
}

func RegisterConrollerInterface(item register.ControllerRegisterItem, w http.ResponseWriter, r *http.Request) interface{} {
	registerBaseController(w, r, &item.Controller)

	return item.Controller
}

func executeControllerDirective(d []string, w http.ResponseWriter, r *http.Request, validation interface{}) {
	item := getControllerItem(d[0])
	container := BuildCustomContainer(item.Modules)

	payload := structToMap(validation)

	err := container.Provide(func() Request {
		return payload
	})

	if err != nil {
		log.Fatal(err)
	}

	cc := RegisterConrollerInterface(item, w, r)
	method := reflect.ValueOf(cc).MethodByName(d[1])

	if err := container.Invoke(method.Interface()); err != nil {
		log.Fatal(err)
	}
}

func structToMap(s interface{}) map[string]interface{} {
	m := make(map[string]interface{})
	j, _ := json.Marshal(s)
	_ = json.Unmarshal(j, &m)

	return m
}

func validateRequest(data interface{}, r *http.Request) error {
	if data != nil {
		if err := tool.DecodeJsonRequest(r, &data); err != nil {
			return err
		}

		if err := tool.ValidateRequest(data); err != nil {
			return err
		}
	}

	return nil
}
