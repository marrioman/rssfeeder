package feed

import (
	"context"
	"fmt"
	"time"

	"github.com/marrioman/rssfeeder/internal/common"

	"github.com/labstack/gommon/log"
	"github.com/mmcdole/gofeed"
)

func AddSource(ctx common.Context, rsslink string) (err error) {

	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(rsslink)
	if err != nil {
		log.Error(err)
		return
	}

	var lastInsertID int64
	sqlStatement := `INSERT INTO feeder.sources (url, sitename) VALUES ($1, $2) ON CONFLICT (sitename) DO NOTHING RETURNING id`
	err = ctx.USERDB.QueryRowx(sqlStatement, rsslink, feed.Title).Scan(&lastInsertID)
	if lastInsertID == 0 {
		err = fmt.Errorf("source already in base")
		return err
	}
	if err != nil {
		log.Error(err)
		return
	}

	sqlStatement2 := `INSERT INTO feeder.content (item_title, link, description, source_id) VALUES ($1, $2, $3, $4)`
	for _, item := range feed.Items {
		_, err = ctx.USERDB.Exec(sqlStatement2, item.Title, item.Link, item.Description, lastInsertID)
		if err != nil {
			log.Error(err)
			return
		}
	}

	fmt.Println(feed.Title, rsslink, lastInsertID)

	return
}

func GetFeed(ctx common.Context) (feed []Item, err error) {

	sqlStatement := `SELECT * FROM feeder.content`
	err = ctx.USERDB.Select(&feed, sqlStatement)
	if err != nil {
		log.Error(err)
		return
	}

	return
}

// TODO: при больших объемах данных добавить индекс:
// CREATE INDEX idx_gin ON item_title USING gin (md5 gin_trgm_ops);
// CREATE INDEX
func SearchByTitle(ctx common.Context, title string) (feed []Item, err error) {

	sqlStatement := `SELECT * FROM feeder.content WHERE item_title ilike $1`
	err = ctx.USERDB.Select(&feed, sqlStatement, "%"+title+"%")
	if err != nil {
		log.Error(err)
		return
	}

	return
}

func GetSourcesUrls(ctx common.Context) (sources []Source, err error) {
	sqlStatement := `SELECT id, url FROM feeder.sources`
	err = ctx.USERDB.Select(&sources, sqlStatement)
	if err != nil {
		log.Error(err)
		return
	}
	return
}

func ParseAndUpdateSources(ctx common.Context, sources []Source) (err error) {
	for _, source := range sources {
		rctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		fp := gofeed.NewParser()
		feed, _ := fp.ParseURLWithContext(source.URL, rctx)
		UpdateContentInDB(ctx, feed, source.ID)
	}
	return
}

func UpdateContentInDB(ctx common.Context, feed *gofeed.Feed, sourceID int) (err error) {
	sqlStatement := `INSERT INTO feeder.content (item_title, link, description, source_id) 
	VALUES ($1, $2, $3, $4)
	ON CONFLICT (item_title) DO NOTHING`
	for _, item := range feed.Items {
		_, err = ctx.USERDB.Exec(sqlStatement, item.Title, item.Link, item.Description, sourceID)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}
