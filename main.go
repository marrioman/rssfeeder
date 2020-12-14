package main

import (
	"log"
	"time"

	"github.com/marrioman/rssfeeder/internal/config"
	"github.com/marrioman/rssfeeder/internal/controller"
	"github.com/marrioman/rssfeeder/internal/database"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/labstack/echo/v4"
)

func main() {
	config.InitConfig()

	_, err := database.InitDatabase()
	if err != nil {
		log.Fatal(err)
	}

	ch, err := migrate.New("file://./migrations", config.C.Database.Feederdb.URL)
	if err != nil {
		log.Fatalln(err)
	}

	err = ch.Up()
	if err != nil && err.Error() != "no changes in database" {
		log.Fatalln(err)
	}

	go func() {
		for {
			controller.FeedUpdater()
			time.Sleep(time.Minute * time.Duration(1))
		}
	}()

	e := echo.New()

	api := e.Group("/api")

	controller.Add(api)

	e.Logger.Fatal(e.Start(config.C.Server.URL))
}
