package main

import (
	"log"

	"github.com/vitalfit/api/config"
	"github.com/vitalfit/api/internal/app"
	"github.com/vitalfit/api/internal/store"
	"github.com/vitalfit/api/pkg/db"

	_ "github.com/lib/pq" // Importa el driver de PostgreSQL
)

//	@title			VitalFit API
//	@description	API for VitalFit, a gym system management
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

// @BasePath					/v1
//
// @securityDefinitions.apikey	ApiKeyAuth
// @in							header
// @name						Authorization
// @description
func main() {
	//initialize config
	config := config.LoadConfig()

	//db gorm connection
	db, err := db.New(config.Db.Dsn, config.Db.MaxOpenConns, config.Db.MaxIdleConns, config.Db.MaxIdleTime)
	if err != nil {
		log.Fatal(err)
	}

	//initialize store
	store := store.NewStorage(db, *config)
	//initialize application

	app := app.BuildApplication(config, store)
	mux := app.Mount()
	if err := app.Run(mux); err != nil {
		log.Fatal(err)
	}

}
