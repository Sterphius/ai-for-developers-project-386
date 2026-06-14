package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
)

type fakeClock struct {
	now time.Time
}

func (clock fakeClock) Now() time.Time {
	return clock.now
}

type fakeRepository struct {
	eventTypes []domain.EventType
	bookings   []domain.Booking
}

func (repository *fakeRepository) CreateEventType(ctx context.Context, eventType domain.EventType) (domain.EventType, error) {
	repository.eventTypes = append(repository.eventTypes, eventType)
	return eventType, nil
}

func (repository *fakeRepository) ListEventTypes(ctx context.Context) ([]domain.EventType, error) {
	return append([]domain.EventType(nil), repository.eventTypes...), nil
}

func (repository *fakeRepository) GetEventType(ctx context.Context, eventTypeID string) (domain.EventType, error) {
	for _, eventType := range repository.eventTypes {
		if eventType.ID == eventTypeID {
			return eventType, nil
		}
	}

	return domain.EventType{}, ErrNotFound
}

func (repository *fakeRepository) CreateBooking(ctx context.Context, booking domain.Booking) error {
	for _, existing := range repository.bookings {
		if booking.Slot.StartAt.Before(existing.Slot.EndAt) && booking.Slot.EndAt.After(existing.Slot.StartAt) {
			return ErrConflict
		}
	}

	repository.bookings = append(repository.bookings, booking)
	return nil
}

func (repository *fakeRepository) ListBookings(ctx context.Context) ([]domain.Booking, error) {
	return append([]domain.Booking(nil), repository.bookings...), nil
}

func (repository *fakeRepository) ListBookingsBetween(ctx context.Context, from time.Time, to time.Time) ([]domain.Booking, error) {
	result := make([]domain.Booking, 0)
	for _, booking := range repository.bookings {
		if booking.Slot.StartAt.Before(to) && booking.Slot.EndAt.After(from) {
			result = append(result, booking)
		}
	}

	return result, nil
}

func TestListAvailableSlotsSkipsOccupiedSlots(t *testing.T) {
	location := time.UTC
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, location)
	repository := &fakeRepository{
		eventTypes: []domain.EventType{
			{ID: "event-1", Name: "Consultation", Description: "Consultation", DurationMinutes: 30},
		},
		bookings: []domain.Booking{
			{
				ID:          "booking-1",
				EventTypeID: "event-1",
				Slot: domain.Slot{
					StartAt: time.Date(2026, 6, 14, 9, 0, 0, 0, location),
					EndAt:   time.Date(2026, 6, 14, 9, 30, 0, 0, location),
				},
			},
		},
	}
	appService := New(repository, Settings{
		OwnerID:             "owner-1",
		Location:            location,
		WindowDays:          14,
		WorkingDayStartHour: 9,
		WorkingDayEndHour:   18,
		SlotStepMinutes:     30,
	})
	appService.clock = fakeClock{now: now}

	slots, err := appService.ListAvailableSlots(context.Background(), "event-1")
	if err != nil {
		t.Fatalf("ListAvailableSlots returned error: %v", err)
	}
	if len(slots) == 0 {
		t.Fatal("expected available slots")
	}
	for _, slot := range slots {
		if slot.StartAt.Equal(repository.bookings[0].Slot.StartAt) {
			t.Fatalf("occupied slot returned as available: %v", slot.StartAt)
		}
	}
}

func TestCreateBookingRejectsPastSlot(t *testing.T) {
	location := time.UTC
	now := time.Date(2026, 6, 14, 10, 0, 0, 0, location)
	repository := &fakeRepository{
		eventTypes: []domain.EventType{
			{ID: "event-1", Name: "Consultation", Description: "Consultation", DurationMinutes: 30},
		},
	}
	appService := New(repository, Settings{
		OwnerID:             "owner-1",
		Location:            location,
		WindowDays:          14,
		WorkingDayStartHour: 9,
		WorkingDayEndHour:   18,
		SlotStepMinutes:     30,
	})
	appService.clock = fakeClock{now: now}

	_, err := appService.CreateBooking(context.Background(), domain.CreateBookingRequest{
		EventTypeID: "event-1",
		SlotStartAt: time.Date(2026, 6, 13, 9, 0, 0, 0, location),
	})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}

func TestCreateEventTypeValidatesInput(t *testing.T) {
	appService := New(&fakeRepository{}, Settings{OwnerID: "owner-1", Location: time.UTC})

	_, err := appService.CreateEventType(context.Background(), domain.CreateEventTypeRequest{})
	if !errors.Is(err, ErrValidation) {
		t.Fatalf("expected validation error, got %v", err)
	}
}
