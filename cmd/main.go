package main

import (
	"github.com/tochamateusz/micro-example/modules/config"
	testservice "github.com/tochamateusz/micro-example/modules/test_service"
	"github.com/tochamateusz/micro/modules/database"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		database.Module(),
		fx.Decorate(config.ProvideDbConfig),
		fx.Invoke(testservice.ProvideTestService),
	).Run()
}
