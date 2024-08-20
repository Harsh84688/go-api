package routes

import (
	"book-manager/middlewares"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(server *gin.Engine) {
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
