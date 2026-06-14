package httpapi

import (
	"context"
	"net/http"
	"os"

	"github.com/Sterphius/ai-for-developers-project-386/internal/domain"
	"github.com/Sterphius/ai-for-developers-project-386/internal/service"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
)

type AppService interface {
	ListPublicEventTypes(ctx context.Context) ([]domain.PublicEventType, error)
	ListAvailableSlots(ctx context.Context, eventTypeID string) ([]domain.Slot, error)
	CreateBooking(ctx context.Context, request domain.CreateBookingRequest) (domain.Booking, error)
	CreateEventType(ctx context.Context, request domain.CreateEventTypeRequest) (domain.EventType, error)
	ListBookings(ctx context.Context) ([]domain.Booking, error)
}

func NewRouter(appService AppService) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	handler := &handler{service: appService}

	router.GET("/healthz", func(context *gin.Context) {
		context.JSON(http.StatusOK, gin.H{"status": "ok"})
	})
	router.GET("/openapi.yaml", serveOpenAPISpec)
	router.GET("/swagger", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/", func(context *gin.Context) {
		context.Redirect(http.StatusMovedPermanently, "/swagger/index.html")
	})
	router.GET("/swagger/index.html", serveSwaggerAsset("/index.html", "text/html; charset=utf-8"))
	router.GET("/swagger/index.css", serveSwaggerAsset("/index.css", "text/css; charset=utf-8"))
	router.GET("/swagger/swagger-initializer.js", serveSwaggerAsset("/swagger-initializer.js", "application/javascript"))
	router.GET("/swagger/swagger-ui.css", serveSwaggerAsset("/swagger-ui.css", "text/css; charset=utf-8"))
	router.GET("/swagger/swagger-ui.css.map", serveSwaggerAsset("/swagger-ui.css.map", "application/json; charset=utf-8"))
	router.GET("/swagger/swagger-ui.js", serveSwaggerAsset("/swagger-ui.js", "application/javascript"))
	router.GET("/swagger/swagger-ui.js.map", serveSwaggerAsset("/swagger-ui.js.map", "application/json; charset=utf-8"))
	router.GET("/swagger/swagger-ui-bundle.js", serveSwaggerAsset("/swagger-ui-bundle.js", "application/javascript"))
	router.GET("/swagger/swagger-ui-bundle.js.map", serveSwaggerAsset("/swagger-ui-bundle.js.map", "application/json; charset=utf-8"))
	router.GET("/swagger/swagger-ui-standalone-preset.js", serveSwaggerAsset("/swagger-ui-standalone-preset.js", "application/javascript"))
	router.GET("/swagger/swagger-ui-standalone-preset.js.map", serveSwaggerAsset("/swagger-ui-standalone-preset.js.map", "application/json; charset=utf-8"))
	router.GET("/swagger/favicon-16x16.png", serveSwaggerAsset("/favicon-16x16.png", "image/png"))
	router.GET("/swagger/favicon-32x32.png", serveSwaggerAsset("/favicon-32x32.png", "image/png"))
	router.GET("/swagger/oauth2-redirect.html", serveSwaggerAsset("/oauth2-redirect.html", "text/html; charset=utf-8"))

	public := router.Group("/public")
	public.GET("/event-types", handler.listPublicEventTypes)
	public.GET("/event-types/:eventTypeId/slots", handler.listAvailableSlots)
	public.POST("/bookings", handler.createBooking)

	owner := router.Group("/owner")
	owner.POST("/event-types", handler.createEventType)
	owner.GET("/bookings", handler.listBookings)

	return router
}

func init() {
	initializer := `window.onload = function() {
  window.ui = SwaggerUIBundle({
    url: "/openapi.yaml",
    dom_id: "#swagger-ui",
    deepLinking: true,
    presets: [
      SwaggerUIBundle.presets.apis,
      SwaggerUIStandalonePreset
    ],
    plugins: [
      SwaggerUIBundle.plugins.DownloadUrl
    ],
    layout: "StandaloneLayout"
  });
};`

	if err := swaggerFiles.WriteFile("/swagger-initializer.js", []byte(initializer), 0o644); err != nil {
		panic(err)
	}
}

func serveSwaggerAsset(assetPath string, contentType string) gin.HandlerFunc {
	return func(context *gin.Context) {
		data, err := swaggerFiles.ReadFile(assetPath)
		if err != nil {
			context.JSON(http.StatusNotFound, domain.ErrorResponse{
				Code:    "not_found",
				Message: "swagger asset not found",
			})
			return
		}

		context.Data(http.StatusOK, contentType, data)
	}
}

func serveOpenAPISpec(context *gin.Context) {
	for _, path := range []string{
		"openapi/openapi.yaml",
		"../openapi/openapi.yaml",
		"../../openapi/openapi.yaml",
	} {
		data, err := os.ReadFile(path)
		if err == nil {
			context.Data(http.StatusOK, "application/yaml; charset=utf-8", data)
			return
		}
	}

	context.JSON(http.StatusNotFound, domain.ErrorResponse{
		Code:    "not_found",
		Message: "openapi specification not found",
	})
}

func mapServiceError(err error) (int, domain.ErrorResponse) {
	switch {
	case err == nil:
		return http.StatusOK, domain.ErrorResponse{}
	case err == service.ErrValidation:
		return http.StatusBadRequest, domain.ErrorResponse{Code: "validation_error", Message: "request validation failed"}
	case err == service.ErrNotFound:
		return http.StatusNotFound, domain.ErrorResponse{Code: "not_found", Message: "resource not found"}
	case err == service.ErrConflict:
		return http.StatusConflict, domain.ErrorResponse{Code: "conflict", Message: "slot is already occupied"}
	default:
		return http.StatusInternalServerError, domain.ErrorResponse{Code: "internal_error", Message: "internal server error"}
	}
}
