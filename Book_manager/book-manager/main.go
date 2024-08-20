package main

import (
	"book-manager/cache"
	"book-manager/db"
	"book-manager/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	cache.InitCache()
	server := gin.Default()
	routes.RegisterRoutes(server)
	server.Run(":8080")
}
