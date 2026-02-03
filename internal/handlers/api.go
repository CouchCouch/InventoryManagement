package handlers

import (
	"inventory/internal/db"
	"net/http"
	"strings"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type APIHandler struct {
	db *db.DB
}

func Handle(r *gin.Engine, db *db.DB) {
	APIHandlerInstance := &APIHandler{db: db}
	api := r.Group("/api")
	{
		itemsAPI := api.Group("/items")
		{
			itemsAPI.GET("", APIHandlerInstance.GetItemsHandler)
			itemsAPI.POST("", APIHandlerInstance.AddItemHandler)
		}
		checkoutsAPI := api.Group("/checkouts")
		{
			checkoutsAPI.GET("", APIHandlerInstance.GetCheckoutsHandler)
			checkoutsAPI.POST("", APIHandlerInstance.CreateCheckoutHandler)
		}
	}
	r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			index, err := static.LocalFile("./web/dist", true).Open("index.html")
			if err != nil {
				log.Error(err)
			}
			defer index.Close()
			stat, _ := index.Stat()
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), index)
		}
	})
}
