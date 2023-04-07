package register

type Route struct {
	Name        string
	Path        string
	Action      string
	Method      string
	Description string
	Validation  interface{}
	Middleware  []Middleware
}

type Group struct {
	Name       string
	Prefix     string
	Routes     []Route
	Middleware []Middleware
}

type HTTPRouter struct {
	Route  []Route
	Groups []Group
}
