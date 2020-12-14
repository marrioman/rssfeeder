package database

import (
	"log"

	"github.com/marrioman/rssfeeder/internal/config"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

var DB *sqlx.DB

func InitDatabase() (*sqlx.DB, error) {
	var err error
	DB, err = sqlx.Open(config.C.Database.Feederdb.Dialect, config.C.Database.Feederdb.URL)
	if err != nil {
		log.Fatal(err)
	}
	DB.SetMaxIdleConns(config.C.Database.Feederdb.Poolsizemin)
	DB.SetMaxOpenConns(config.C.Database.Feederdb.Poolsizemax)
	err = DB.Ping()
	if err != nil {
		log.Fatal(err)
	}

	return DB, err
}
