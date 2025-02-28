package main

import (
	"b8boost/backend/config"
	"b8boost/backend/internal/infra"
)

//	@title			B8boost API
//	@version		1.0
//	@BasePath		/

// @securityDefinitions.apikey	BearerAuth
// @in							header
// @name						Authorization
func main() {
	config, err := config.NewLoadConfig()
	if err != nil {
		panic(err)
	}

	app := infra.Config(config).Database().JWT().Ldap().LLM().TgBot().Cron().Serve()

	app.Start()
}
