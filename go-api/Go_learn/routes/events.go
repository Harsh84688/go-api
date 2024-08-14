package routes

import (
	"Go_learn/Go_learn/models"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func getEvents(context *gin.Context) {
	events, err := models.GetAllEvents()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not get events"})
		return
	}
	context.JSON(http.StatusOK, events)
}

func getEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse eventID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event"})
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create event"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "event created", "event": event})
}

func updateEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse eventID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event"})
		return
	}

	if event.UserId != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized to update this event"})
		return
	}

	var updatedEvent models.Event
	err = context.ShouldBindJSON(&updatedEvent)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse event"})
		return
	}

	updatedEvent.ID = eventId

	err = updatedEvent.Update()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not update event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "event updated", "event": updatedEvent})
}

func deleteEvent(context *gin.Context) {
	eventId, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"message": "could not parse eventID"})
		return
	}

	event, err := models.GetEventByID(eventId)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not fetch event"})
		return
	}

	if event.UserId != context.GetInt64("userId") {
		context.JSON(http.StatusUnauthorized, gin.H{"message": "unauthorized to delete this event"})
		return
	}

	err = event.Delete()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not delete event"})
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "event deleted"})
}

func createManyHelper(context *gin.Context, event *models.Event, eventNum string, errorChannel chan string) {
	if eventNum == "event1" {
		errorChannel <- eventNum
		return
	}

	event.UserId = context.GetInt64("userId")

	err := event.Save()

	if err != nil {
		errorChannel <- eventNum
		return
	}

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
		go createManyHelper(context, &event, eventNum, errorChannel)
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
		context.JSON(http.StatusInternalServerError, gin.H{"message": "could not create these events", "eventNums": codes})
		return
	}
	context.JSON(http.StatusCreated, gin.H{"message": "events created", "events": events})
}
