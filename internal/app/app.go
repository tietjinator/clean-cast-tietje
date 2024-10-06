package app

import (
	setup "ikoyhn/podcast-sponsorblock/internal/database"
	"ikoyhn/podcast-sponsorblock/internal/services"
	"net/http"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type Env struct {
	db *gorm.DB
}

func Start() {
	db, err := setup.ConnectDatabase()
	if err != nil {
		log.Panic(err)
	}
	env := &Env{db: db}

	e := echo.New()

	e.Use(middleware.Logger())

	env.registerRoutes(e)
}

func (env *Env) registerRoutes(e *echo.Echo) {
	e.GET("/rss/:youtubePlaylistId", func(c echo.Context) error {
		data := services.BuildRssFeed(env.db, c, c.Param("youtubePlaylistId"))
		return c.JSON(http.StatusOK, data)
	})
	e.GET("/media/:youtubeVideoId", func(c echo.Context) error {
		return services.GetYoutubeVideo(c.Param("youtubeVideoId"), c)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("HOST")

	log.Info("Starting server on " + host + ": " + port)
	e.Logger.Fatal(e.Start(host + ":" + port))

}
