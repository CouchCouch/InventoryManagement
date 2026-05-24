package handlers

import (
	"log/slog"
	"net/http"
	"strings"

	"inventory/internal/auth"
	"inventory/internal/db"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
)

const htmlResponse string = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Error</title>
    <style>
        body {
            font-family: Arial, sans-serif;
            background-color: #f8d7da;
            color: #721c24;
            text-align: center;
            padding: 50px;
        }
        h1 {
            font-size: 36px;
        }
        p {
            font-size: 20px;
        }
        a {
            color: #0056b3;
            text-decoration: none;
        }
        a:hover {
            text-decoration: underline;
        }
    </style>
</head>
<body>
    <h1>Oops! Something went wrong.</h1>
    <p>We're sorry, but an error has occurred while processing your request.</p>
    <p>Please try again later or <a href="/">return to the homepage</a>.</p>
</body>
</html>
`

type APIHandler struct {
	db   *db.DB
	auth *auth.AuthService
}

func Handle(r *gin.Engine, db *db.DB, auth *auth.AuthService) {
	// r.Use(RequestLoggingMiddleware())

	APIHandlerInstance := &APIHandler{db: db, auth: auth}
	api := r.Group("/api")
	{
		itemsAPI := api.Group("/items")
		{
			itemsAPI.GET("", APIHandlerInstance.GetItemsHandler)
			itemsAPI.GET("/status", APIHandlerInstance.GetItemsStatusHandler)
			itemsAPI.GET("/types", APIHandlerInstance.GetItemsTypes)
			itemsAPI.POST("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.AddItemHandler)
			itemsAPI.DELETE("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.DeleteItemHandler)
		}
		checkoutsAPI := api.Group("/checkouts")
		{
			checkoutsAPI.GET("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.GetCheckoutsHandler)
			checkoutsAPI.POST("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.CreateCheckoutHandler)
			// checkoutsAPI.PUT("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.ReturnCheckoutHandler)
			// checkoutsAPI.PATCH("", APIHandlerInstance.AuthMiddleware(), )
		}
		authAPI := api.Group("/auth")
		{
			authAPI.POST("/login", APIHandlerInstance.LoginHandler)
			authAPI.POST("/logout", APIHandlerInstance.LogoutHandler)
		}
		userAPI := api.Group("/users")
		{
			userAPI.GET("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.GetUserHandler)
			userAPI.POST("", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.CreateUserHandler)
			userAPI.GET("/me", APIHandlerInstance.AuthMiddleware(), APIHandlerInstance.GetMeHandler)
		}
	}
	r.Use(static.Serve("/", static.LocalFile("./web/dist", true)))
	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.URL.Path, "/api") {
			index, err := static.LocalFile("./web/dist", true).Open("index.html")
			if err != nil {
				slog.Error("Error loading file", "error", err, "file", "index.html")
				c.HTML(500, htmlResponse, gin.H{})
				return
			}

			//nolint:errcheck
			defer index.Close()

			stat, err := index.Stat()
			if err != nil {
				slog.Error("Error getting file info", "error", err, "file", "index.html")
				c.HTML(500, htmlResponse, gin.H{})
				return
			}
			http.ServeContent(c.Writer, c.Request, "index.html", stat.ModTime(), index)
		}
	})
}
