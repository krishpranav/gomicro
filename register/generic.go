package register

type ServiceRegister []interface{}

type ModelRegister []interface{}

type ControllerRegister []ControllerRegisterItem

type ControllerRegisterItem struct {
	Controller interface{}
	Modules    []DIModule
}
