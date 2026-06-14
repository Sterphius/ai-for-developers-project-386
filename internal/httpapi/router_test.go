package httpapi

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
	"github.com/gin-gonic/gin"
)

type stubService struct{}

func (stubService) ListPublicEventTypes(ctx context.Context) ([]domain.PublicEventType, error) {
	return []domain.PublicEventType{{ID: "event-1", Name: "Consultation", Description: "Consultation", DurationMinutes: 30}}, nil
}

func (stubService) ListAvailableSlots(ctx context.Context, eventTypeID string) ([]domain.Slot, error) {
	return []domain.Slot{{StartAt: time.Now().UTC(), EndAt: time.Now().UTC().Add(30 * time.Minute), IsFree: true}}, nil
}

func (stubService) CreateBooking(ctx context.Context, request domain.CreateBookingRequest) (domain.Booking, error) {
	return domain.Booking{ID: "booking-1", EventTypeID: request.EventTypeID}, nil
}

func (stubService) CreateEventType(ctx context.Context, request domain.CreateEventTypeRequest) (domain.EventType, error) {
	return domain.EventType{ID: request.ID, Name: request.Name}, nil
}

func (stubService) ListBookings(ctx context.Context) ([]domain.Booking, error) {
	return []domain.Booking{{ID: "booking-1"}}, nil
}

func TestHealthEndpoint(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(stubService{})

	request, _ := http.NewRequest(http.MethodGet, "/healthz", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", responseRecorder.Code)
	}
}

func TestOpenAPIDocumentIsServed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(stubService{})

	request, _ := http.NewRequest(http.MethodGet, "/openapi.yaml", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", responseRecorder.Code)
	}
	if responseRecorder.Body.Len() == 0 {
		t.Fatal("expected openapi document body")
	}
}

func TestSwaggerUIIsServed(t *testing.T) {
	gin.SetMode(gin.TestMode)
	router := NewRouter(stubService{})

	request, _ := http.NewRequest(http.MethodGet, "/swagger/index.html", nil)
	responseRecorder := httptest.NewRecorder()
	router.ServeHTTP(responseRecorder, request)

	if responseRecorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", responseRecorder.Code)
	}
}
