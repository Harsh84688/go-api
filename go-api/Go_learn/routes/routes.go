package routes

import (
	"Go_learn/Go_learn/middlewares"

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
	authenticated.POST("/events/:id/register", registerForEvent)
	authenticated.DELETE("/events/:id/register", unregisterForEvent)

	server.POST("signup", signup)
	server.POST("login", login)
}
