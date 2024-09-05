package middlewares

import (
	"book-manager/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func Auth(context *gin.Context) {
	// token := context.Request.Header.Get("Authorization")
	// if token == "" {
	// 	context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
	// 	return
	// }

	token, err := context.Cookie("token")
	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "missing token"})
		return
	}

	userId, err := utils.ValidateToken(token)
	log.Println(userId)

	if err != nil {
		context.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"message": "unauthorized"})
		return
	}

	context.Set("userId", userId)

	context.Next()
}
