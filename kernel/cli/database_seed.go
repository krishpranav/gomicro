package cli

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/krishpranav/gomicro/register"

	"github.com/jinzhu/gorm"
)

type Seeder struct {
	register.Command
}

func (c *Seeder) Register() {
	c.Signature = "database:seed <name>"
	c.Description = "Execute database seeder"
}

func (c *Seeder) Run(db *gorm.DB, models register.ModelRegister) {
	fmt.Println("Execute database seeding...")
	if len(c.Args) > 0 {
		extractSpecificModel(c.Args, &models)
	}

	seed(models, db)
}

func extractSpecificModel(name string, models *register.ModelRegister) {
	var newModels register.ModelRegister

	for _, m := range *models {
		modelName := reflect.TypeOf(m).Name()

		if strings.EqualFold(name, modelName) {
			newModels = append(newModels, m)
			break
		}
	}

	*models = newModels
}

func seed(models []interface{}, db *gorm.DB) {
	for _, m := range models {
		fmt.Printf("\nCreating items for model %s...\n", reflect.TypeOf(m).Name())
		v := reflect.ValueOf(m)
		method := v.MethodByName("Seed")
		method.Call([]reflect.Value{reflect.ValueOf(db)})

		fmt.Printf("Success!\n")
	}

	fmt.Println("Seeding complete!")
}
