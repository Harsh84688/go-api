package routes

import (
	"book-manager/models"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get books"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse bookID"})
		return
	}

	event, err := models.GetEventByID(eventId, context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch book"})
		return
	}
	context.JSON(http.StatusOK, event)
}

func createEvent(context *gin.Context) {

	var event models.Event
	err := context.ShouldBindJSON(&event)

	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse"})
		return
	}

	event.UserId = context.GetInt64("userId")

	err = event.Save()

	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create book"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "book created", "book": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse bookID"})
		return
	}

	event, err := models.GetEventByID(eventId, context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch book"})
		return
	}

	if event.UserId != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized to update this book"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse book"})
		return
	}

	updatedEvent.ID = eventId
	updatedEvent.UserId = event.UserId

	err = updatedEvent.Update(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update book"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "book updated", "event": updatedEvent})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse bookID"})
		return
	}

	event, err := models.GetEventByID(eventId, context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch book"})
		return
	}

	if event.UserId != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized to delete this book"})
		return
	}

	err = event.Delete(context.Param("id"))
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete book"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "book deleted"})
}

func createManyHelper(context *gin.Context, event *models.Event, eventNum string, errorChannel chan string, events *map[string]models.Event) {

	event.UserId = context.GetInt64("userId")

	err := event.Save()

	if err != nil {
		log.Println(err)
		errorChannel <- eventNum
		return
	}
	delete(*events, eventNum)
	(*events)[eventNum] = *event

	errorChannel <- ""
}

func createManyEvents(context *gin.Context) {

	var events map[string]models.Event

	err := context.ShouldBindJSON(&events)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse"})
		return
	}

	var errorChannel = make(chan string, len(events))

	for eventNum, event := range events {
		go createManyHelper(context, &event, eventNum, errorChannel, &events)
	}

	var codes []string
	var temp string
	for i := 0; i < len(events); i++ {
		temp = <-errorChannel
		if temp != "" {
			codes = append(codes, temp)
		}
	}

	if len(codes) > 0 {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create these books", "bookNums": codes})
		return
	}
	log.Println(events)
	context.JSON(http.StatusCreated, gin.H{"message": "books created", "books": events})
}
