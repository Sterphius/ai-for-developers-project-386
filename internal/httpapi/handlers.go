package httpapi

import (
	"net/http"
	"time"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
	"github.com/gin-gonic/gin"
)

type handler struct {
	service AppService
}

type createBookingInput struct {
	EventTypeID string `json:"eventTypeId"`
	SlotStartAt string `json:"slotStartAt"`
}

func (handler *handler) listPublicEventTypes(context *gin.Context) {
	eventTypes, err := handler.service.ListPublicEventTypes(context.Request.Context())
	if err != nil {
		statusCode, response := mapServiceError(err)
		context.JSON(statusCode, response)
		return
	}

	context.JSON(http.StatusOK, eventTypes)
}

func (handler *handler) listAvailableSlots(context *gin.Context) {
	eventTypeID := context.Param("eventTypeId")
	slots, err := handler.service.ListAvailableSlots(context.Request.Context(), eventTypeID)
	if err != nil {
		statusCode, response := mapServiceError(err)
		context.JSON(statusCode, response)
		return
	}

	context.JSON(http.StatusOK, slots)
}

func (handler *handler) createBooking(context *gin.Context) {
	var input createBookingInput
	if err := context.ShouldBindJSON(&input); err != nil {
		context.JSON(http.StatusBadRequest, domain.ErrorResponse{Code: "validation_error", Message: "invalid request body"})
		return
	}

	slotStartAt, err := time.Parse(time.RFC3339, input.SlotStartAt)
	if err != nil {
		context.JSON(http.StatusBadRequest, domain.ErrorResponse{Code: "validation_error", Message: "slotStartAt must be RFC3339"})
		return
	}

	booking, err := handler.service.CreateBooking(context.Request.Context(), domain.CreateBookingRequest{
		EventTypeID: input.EventTypeID,
		SlotStartAt: slotStartAt,
	})
	if err != nil {
		statusCode, response := mapServiceError(err)
		context.JSON(statusCode, response)
		return
	}

	context.JSON(http.StatusCreated, booking)
}

func (handler *handler) createEventType(context *gin.Context) {
	var request domain.CreateEventTypeRequest
	if err := context.ShouldBindJSON(&request); err != nil {
		context.JSON(http.StatusBadRequest, domain.ErrorResponse{Code: "validation_error", Message: "invalid request body"})
		return
	}

	eventType, err := handler.service.CreateEventType(context.Request.Context(), request)
	if err != nil {
		statusCode, response := mapServiceError(err)
		context.JSON(statusCode, response)
		return
	}

	context.JSON(http.StatusCreated, eventType)
}

func (handler *handler) listBookings(context *gin.Context) {
	bookings, err := handler.service.ListBookings(context.Request.Context())
	if err != nil {
		statusCode, response := mapServiceError(err)
		context.JSON(statusCode, response)
		return
	}

	context.JSON(http.StatusOK, bookings)
}
