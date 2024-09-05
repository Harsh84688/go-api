package routes

import (
	"book-manager/middlewares"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
	viewspath := filepath.Join("..", "app", "views")
	server.Static("/static", viewspath)

	// Serve index.html at the root path
	server.GET("/", func(c *gin.Context) {
		c.File(filepath.Join("views", "index.html"))
	})

	// Serve other static files (e.g., dashboard.html) if needed
	server.GET("/dashboard", func(c *gin.Context) {
		c.File(filepath.Join("views", "dashboard.html"))
	})

	server.GET("/events", getEvents)
	server.GET("/events/:id", getEvent)

	authenticated := server.Group("/")
	authenticated.Use(middlewares.Auth, middlewares.RateLimitMiddleware)

	authenticated.POST("/events", createEvent)
	authenticated.POST("events/many", createManyEvents)
	authenticated.PUT("/events/:id", updateEvent)
	authenticated.DELETE("/events/:id", deleteEvent)

	server.POST("signup", signup)
	server.POST("login", login)
}
