package controller

import (
	"log"
	"net/http"
	"github.com/marrioman/rssfeeder/internal/common"
	"github.com/marrioman/rssfeeder/internal/feed"

	"github.com/labstack/echo/v4"
)

func Add(api *echo.Group) {

	api.POST("/addSource", addSource)
	api.GET("/feed", getFeed)
	api.GET("/search", searchByTitle)

	return
}

func addSource(c echo.Context) (err error) {
	var (
		req struct {
			Rss string `json:"rss"`
		}
		ctx = common.NewContext(c)
	)

	err = c.Bind(&req)
	if err != nil {
		return err
	}

	err = feed.AddSource(ctx, req.Rss)
	if err != nil {
		return c.JSON(http.StatusOK, echo.Map{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": "success",
	})
}

func getFeed(c echo.Context) (err error) {

	var ctx = common.NewContext(c)

	response, err := feed.GetFeed(ctx)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": response,
	})
}

func searchByTitle(c echo.Context) (err error) {

	var ctx = common.NewContext(c)

	response, err := feed.SearchByTitle(ctx, c.QueryParam("title"))
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, echo.Map{
		"data": response,
	})
}

func FeedUpdater() (err error) {

	var (
		c   echo.Context
		ctx = common.NewContext(c)
	)

	sources, err := feed.GetSourcesUrls(ctx)
	if err != nil {
		return
	}

	if len(sources) == 0 {
		log.Println("0 sources in database")
		return
	}

	err = feed.ParseAndUpdateSources(ctx, sources)
	if err != nil {
		return
	}
	return

}
