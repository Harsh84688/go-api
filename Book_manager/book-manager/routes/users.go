package routes

import (
	"book-manager/models"
	"book-manager/utils"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func signup(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse"})
		return
	}

	err = user.Save()

	if err != nil {
		log.Print(err)
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create user"})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "user created", "user": user})
}

func login(context *gin.Context) {
	var user models.User

	err := context.ShouldBindJSON(&user)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse"})
		return
	}

	err = user.Authenticate()

	if err != nil {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "could not authenticate user"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not generate token"})
		return
	}

	context.SetCookie("token", token, 7200, "/", "localhost", false, true)
	context.JSON(http.StatusOK, gin.H{"message": "authorization cookie set for 2 hours"})
}
