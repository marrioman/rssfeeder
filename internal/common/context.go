package common

import (
	"github.com/marrioman/rssfeeder/internal/database"

	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"
)

type Context struct {
	USERDB *sqlx.DB
}

func NewContext(c echo.Context) (ctx Context) {
	ctx.USERDB, _ = database.InitDatabase()
	return
}
