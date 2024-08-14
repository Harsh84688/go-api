package main

import (
	"Go_learn/Go_learn/db"
	"Go_learn/Go_learn/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	db.InitDB()
	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8090")
}
