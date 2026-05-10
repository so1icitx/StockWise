package httpapi

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/so1icitx/StockWise/internal/application"
	"github.com/so1icitx/StockWise/internal/config"
)

func TestHealthEndpoint(t *testing.T) {
	router := NewRouter(config.Config{
		AppName:    "StockWise",
		AppEnv:     "test",
		ServerHost: "127.0.0.1",
		ServerPort: "0",
	}, application.Services{})

	recorder := httptest.NewRecorder()
	request := httptest.NewRequest(http.MethodGet, "/health", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK {
		t.Fatalf("expected status %d, got %d", http.StatusOK, recorder.Code)
	}
}
