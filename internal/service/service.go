package service

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
)

type Clock interface {
	Now() time.Time
}

type clock struct{}

func (clock) Now() time.Time {
	return time.Now()
}

type Repository interface {
	CreateEventType(ctx context.Context, eventType domain.EventType) (domain.EventType, error)
	ListEventTypes(ctx context.Context) ([]domain.EventType, error)
	GetEventType(ctx context.Context, eventTypeID string) (domain.EventType, error)
	CreateBooking(ctx context.Context, booking domain.Booking) error
	ListBookings(ctx context.Context) ([]domain.Booking, error)
	ListBookingsBetween(ctx context.Context, from time.Time, to time.Time) ([]domain.Booking, error)
}

type Settings struct {
	OwnerID             string
	Location            *time.Location
	WindowDays          int
	WorkingDayStartHour int
	WorkingDayEndHour   int
	SlotStepMinutes     int
}

type Service struct {
	repository Repository
	clock      Clock
	settings   Settings
}

func New(repository Repository, settings Settings) *Service {
	if settings.WindowDays <= 0 {
		settings.WindowDays = 14
	}
	if settings.WorkingDayStartHour <= 0 {
		settings.WorkingDayStartHour = 9
	}
	if settings.WorkingDayEndHour <= 0 {
		settings.WorkingDayEndHour = 18
	}
	if settings.SlotStepMinutes <= 0 {
		settings.SlotStepMinutes = 30
	}
	if settings.Location == nil {
		settings.Location = time.UTC
	}

	return &Service{
		repository: repository,
		clock:      clock{},
		settings:   settings,
	}
}

func (service *Service) ListPublicEventTypes(ctx context.Context) ([]domain.PublicEventType, error) {
	eventTypes, err := service.repository.ListEventTypes(ctx)
	if err != nil {
		return nil, err
	}

	publicTypes := make([]domain.PublicEventType, 0, len(eventTypes))
	for _, eventType := range eventTypes {
		publicTypes = append(publicTypes, domain.PublicEventType{
			ID:              eventType.ID,
			Name:            eventType.Name,
			Description:     eventType.Description,
			DurationMinutes: eventType.DurationMinutes,
		})
	}

	return publicTypes, nil
}

func (service *Service) CreateEventType(ctx context.Context, request domain.CreateEventTypeRequest) (domain.EventType, error) {
	if err := validateCreateEventTypeRequest(request); err != nil {
		return domain.EventType{}, err
	}

	eventType := domain.EventType{
		ID:              request.ID,
		OwnerID:         service.settings.OwnerID,
		Name:            request.Name,
		Description:     request.Description,
		DurationMinutes: request.DurationMinutes,
	}

	return service.repository.CreateEventType(ctx, eventType)
}

func (service *Service) ListAvailableSlots(ctx context.Context, eventTypeID string) ([]domain.Slot, error) {
	eventType, err := service.repository.GetEventType(ctx, eventTypeID)
	if err != nil {
		return nil, err
	}

	windowStart, windowEnd := service.windowBounds()
	now := service.clock.Now().In(service.settings.Location)
	bookings, err := service.repository.ListBookingsBetween(ctx, windowStart, windowEnd)
	if err != nil {
		return nil, err
	}

	return generateAvailableSlots(windowStart, windowEnd, now, eventType.DurationMinutes, service.settings, bookings), nil
}

func (service *Service) CreateBooking(ctx context.Context, request domain.CreateBookingRequest) (domain.Booking, error) {
	if strings.TrimSpace(request.EventTypeID) == "" {
		return domain.Booking{}, ErrValidation
	}
	if request.SlotStartAt.IsZero() {
		return domain.Booking{}, ErrValidation
	}

	eventType, err := service.repository.GetEventType(ctx, request.EventTypeID)
	if err != nil {
		return domain.Booking{}, err
	}

	windowStart, windowEnd := service.windowBounds()
	now := service.clock.Now().In(service.settings.Location)
	requestStart := request.SlotStartAt.In(service.settings.Location)
	requestEnd := requestStart.Add(time.Duration(eventType.DurationMinutes) * time.Minute)
	if requestStart.Before(windowStart) || requestStart.Before(now) || requestEnd.After(windowEnd) {
		return domain.Booking{}, ErrValidation
	}

	if !isValidSlotStart(requestStart, service.settings) {
		return domain.Booking{}, ErrValidation
	}
	if !fitsWorkingWindow(requestStart, requestEnd, service.settings) {
		return domain.Booking{}, ErrValidation
	}

	booking := domain.Booking{
		ID:          newID(),
		EventTypeID: request.EventTypeID,
		Slot: domain.Slot{
			StartAt: requestStart,
			EndAt:   requestEnd,
			IsFree:  false,
		},
		CreatedAt: service.clock.Now().In(service.settings.Location),
	}

	if err := service.repository.CreateBooking(ctx, booking); err != nil {
		return domain.Booking{}, err
	}

	return booking, nil
}

func (service *Service) ListBookings(ctx context.Context) ([]domain.Booking, error) {
	return service.repository.ListBookings(ctx)
}

func (service *Service) windowBounds() (time.Time, time.Time) {
	now := service.clock.Now().In(service.settings.Location)
	dayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, service.settings.Location)
	return dayStart, dayStart.AddDate(0, 0, service.settings.WindowDays)
}

func validateCreateEventTypeRequest(request domain.CreateEventTypeRequest) error {
	if strings.TrimSpace(request.ID) == "" {
		return ErrValidation
	}
	if strings.TrimSpace(request.Name) == "" {
		return ErrValidation
	}
	if strings.TrimSpace(request.Description) == "" {
		return ErrValidation
	}
	if request.DurationMinutes <= 0 {
		return ErrValidation
	}

	return nil
}

func generateAvailableSlots(windowStart, windowEnd, now time.Time, durationMinutes int, settings Settings, bookings []domain.Booking) []domain.Slot {
	if durationMinutes <= 0 {
		return nil
	}

	sort.Slice(bookings, func(left, right int) bool {
		return bookings[left].Slot.StartAt.Before(bookings[right].Slot.StartAt)
	})

	candidates := make([]domain.Slot, 0)
	duration := time.Duration(durationMinutes) * time.Minute
	step := time.Duration(settings.SlotStepMinutes) * time.Minute

	for currentDay := windowStart; currentDay.Before(windowEnd); currentDay = currentDay.AddDate(0, 0, 1) {
		dayStart := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), settings.WorkingDayStartHour, 0, 0, 0, settings.Location)
		dayEnd := time.Date(currentDay.Year(), currentDay.Month(), currentDay.Day(), settings.WorkingDayEndHour, 0, 0, 0, settings.Location)
		candidateStart := dayStart
		if sameCalendarDay(currentDay, now) {
			candidateStart = alignForward(maxTime(dayStart, now), time.Duration(settings.SlotStepMinutes)*time.Minute)
		}

		for ; !candidateStart.Add(duration).After(dayEnd); candidateStart = candidateStart.Add(step) {
			candidateEnd := candidateStart.Add(duration)
			if overlapsAny(candidateStart, candidateEnd, bookings) {
				continue
			}

			candidates = append(candidates, domain.Slot{
				StartAt: candidateStart,
				EndAt:   candidateEnd,
				IsFree:  true,
			})
		}
	}

	return candidates
}

func overlapsAny(candidateStart, candidateEnd time.Time, bookings []domain.Booking) bool {
	for _, booking := range bookings {
		if candidateStart.Before(booking.Slot.EndAt) && candidateEnd.After(booking.Slot.StartAt) {
			return true
		}
	}

	return false
}

func isValidSlotStart(start time.Time, settings Settings) bool {
	if start.Minute()%settings.SlotStepMinutes != 0 || start.Second() != 0 || start.Nanosecond() != 0 {
		return false
	}

	return true
}

func fitsWorkingWindow(start, end time.Time, settings Settings) bool {
	dayStart := time.Date(start.Year(), start.Month(), start.Day(), settings.WorkingDayStartHour, 0, 0, 0, settings.Location)
	dayEnd := time.Date(start.Year(), start.Month(), start.Day(), settings.WorkingDayEndHour, 0, 0, 0, settings.Location)

	return !start.Before(dayStart) && !end.After(dayEnd)
}

func sameCalendarDay(left, right time.Time) bool {
	return left.Year() == right.Year() && left.YearDay() == right.YearDay()
}

func maxTime(left, right time.Time) time.Time {
	if left.After(right) {
		return left
	}

	return right
}

func alignForward(value time.Time, step time.Duration) time.Time {
	if step <= 0 {
		return value
	}

	stepMinutes := int(step / time.Minute)
	if stepMinutes <= 0 {
		return value
	}

	totalMinutes := value.Hour()*60 + value.Minute()
	remainder := totalMinutes % stepMinutes
	if remainder == 0 && value.Second() == 0 && value.Nanosecond() == 0 {
		return time.Date(value.Year(), value.Month(), value.Day(), value.Hour(), value.Minute(), 0, 0, value.Location())
	}

	alignedMinutes := totalMinutes + stepMinutes - remainder
	return time.Date(
		value.Year(),
		value.Month(),
		value.Day(),
		alignedMinutes/60,
		alignedMinutes%60,
		0,
		0,
		value.Location(),
	)
}

func newID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		panic(fmt.Sprintf("generate id: %v", err))
	}

	return hex.EncodeToString(bytes)
}
