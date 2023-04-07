package register

type CommandRegister map[string]interface{}

type Command struct {
	Signature   string
	Description string
	Args        string
}
