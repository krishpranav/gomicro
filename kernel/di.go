package kernel

import (
	"log"

	"github.com/krishpranav/gomicro/register"
	"go.uber.org/dig"
)

func BuildCustomContainer(modules []register.DIModule) *dig.Container {
	container := dig.New()

	for _, m := range modules {
		for _, p := range m.Provides {
			if err := container.Provide(p); err != nil {
				log.Fatal(err)
			}
		}
	}

	return container
}

func BuildCommandContainer() *dig.Container {
	container := dig.New()

	for _, s := range CommandServices {
		if err := container.Provide(s); err != nil {
			log.Fatal(err)
		}
	}
	injectBasicEntities(container)

	return container
}

func injectBasicEntities(sc *dig.Container) {
	_ = sc.Provide(func() register.ControllerRegister {
		return Controllers
	})

	_ = sc.Provide(func() register.CommandRegister {
		return Commands
	})

	_ = sc.Provide(func() register.ModelRegister {
		return Models
	})

	_ = sc.Provide(func() []register.HTTPRouter {
		return Router
	})
}
