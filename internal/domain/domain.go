package domain

import "time"

type Owner struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type EventType struct {
	ID              string `json:"id"`
	OwnerID         string `json:"ownerId"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
}

type PublicEventType struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
}

type Slot struct {
	StartAt time.Time `json:"startAt"`
	EndAt   time.Time `json:"endAt"`
	IsFree  bool      `json:"isFree"`
}

type Booking struct {
	ID          string    `json:"id"`
	EventTypeID string    `json:"eventTypeId"`
	Slot        Slot      `json:"slot"`
	CreatedAt   time.Time `json:"createdAt"`
}

type BookingWithEventType struct {
	Booking   Booking         `json:"booking"`
	EventType PublicEventType `json:"eventType"`
}

type CreateEventTypeRequest struct {
	ID              string `json:"id"`
	Name            string `json:"name"`
	Description     string `json:"description"`
	DurationMinutes int    `json:"durationMinutes"`
}

type CreateBookingRequest struct {
	EventTypeID string    `json:"eventTypeId"`
	SlotStartAt time.Time `json:"slotStartAt"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
