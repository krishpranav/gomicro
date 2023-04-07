package kernel

import (
	"github.com/krishpranav/gomicro/kernel/cli"
	"github.com/krishpranav/gomicro/register"
)

var (
	Commands = register.CommandRegister{
		"database:seed":      &cli.Seeder{},
		"show:commands":      &cli.ShowCommands{},
		"cmd:create":         &cli.CmdCreate{},
		"controller:create":  &cli.ControllerCreate{},
		"generate:key":       &cli.GenerateKey{},
		"middleware:create":  &cli.MiddlewareCreate{},
		"migration:create":   &cli.MigrationCreate{},
		"migration:rollback": &cli.MigrateRollback{},
		"migration:up":       &cli.MigrationUp{},
		"model:create":       &cli.ModelCreate{},
		"router:show":        &cli.RouterShow{},
		"service:create":     &cli.ServiceCreate{},
		"update":             &cli.UpdateAlfred{},
	}
	CommandServices = register.ServiceRegister{}
	Models          = register.ModelRegister{}
	Controllers     = register.ControllerRegister{}
	Middlewares     = register.MiddlewareRegister{}
	Router          []register.HTTPRouter
)
