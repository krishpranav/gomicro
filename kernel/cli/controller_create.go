package cli

import (
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/krishpranav/gomicro/register"
	"github.com/krishpranav/gomicro/tool"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type ControllerCreate struct {
	register.Command
}

func (c *ControllerCreate) Register() {
	c.Signature = "controller:create <name>"
	c.Description = "Create new controller"
}

func (c *ControllerCreate) Run() {
	fmt.Println("Creating new controller...")
	var _, filename, _, _ = runtime.Caller(0)

	cName := cases.Title(language.Und, cases.NoLower).String(c.Args)
	input, _ := os.ReadFile(filepath.Join(path.Dir(filename), "raw/controller.raw"))

	cContent := strings.ReplaceAll(string(input), "@@TMP@@", cName)
	cFile := fmt.Sprintf("%s/%s.go", tool.GetDynamicPath("app/http/controller"), strings.ToLower(c.Args))
	if err := os.WriteFile(cFile, []byte(cContent), 0755); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("\nSUCCESS: Your %sController has been created at %s", cName, cFile)
	fmt.Printf("\nDO NOT FORGET TO REGISTER IT!")
}
