package http

import (
	"context"
	"api-server/domain/analysis"
	"api-server/domain/mocks"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func setupRouter(t *testing.T) (*gin.Engine, *mocks.MockBuscaCEPAPIClient, *mocks.MockWeatherAPIClient) {
	gin.SetMode(gin.TestMode)
	logger := log.New(os.Stdout, "test-handler - ", log.LstdFlags)

	mockBuscaCEPClient := &mocks.MockBuscaCEPAPIClient{}
	mockWeatherClient := &mocks.MockWeatherAPIClient{}

	analysisService := analysis.NewAnalysisService(mockBuscaCEPClient, mockWeatherClient, logger)
	
	// Manually create the handler struct, similar to what NewHandler does
	handler := &handler{
		analisysService: analysisService,
		log:             logger,
	}

	router := gin.Default()
	// Use the correct route as defined in handler.go
	router.GET("/tempForCep/:cep", handler.RunAnalysis)

	return router, mockBuscaCEPClient, mockWeatherClient
}

func TestHandler_RunAnalysis(t *testing.T) {
	t.Run("should return temperature successfully", func(t *testing.T) {
		router, mockBuscaCEP, mockWeather := setupRouter(t)

		mockBuscaCEP.GetBrasilAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "São Paulo,SP", nil
		}
		mockBuscaCEP.GetViaAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "São Paulo,SP", nil
		}
		mockWeather.GetHGWeatherAPIFunc = func(ctx context.Context, city string) (int, error) {
			assert.Equal(t, "São Paulo,SP", city)
			return 25, nil
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tempForCep/01001000", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var response map[string]float64
		err := json.Unmarshal(w.Body.Bytes(), &response)
		assert.NoError(t, err)
		assert.Equal(t, float64(25), response["temp_C"])
		assert.Equal(t, float64(77), response["temp_F"])
		assert.InDelta(t, 298.15, response["temp_K"], 0.01)
	})

	t.Run("should return 422 for invalid cep", func(t *testing.T) {
		router, _, _ := setupRouter(t)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tempForCep/123", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Code)
		assert.JSONEq(t, `{"error":"invalid zipcode"}`, w.Body.String())
	})

	t.Run("should return 404 when cep is not found", func(t *testing.T) {
		router, mockBuscaCEP, _ := setupRouter(t)

		mockBuscaCEP.GetBrasilAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "", errors.New("not found")
		}
		mockBuscaCEP.GetViaAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "", errors.New("not found")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tempForCep/99999999", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
		assert.JSONEq(t, `{"error":"can not find zipcode"}`, w.Body.String())
	})

	t.Run("should return 500 when weather api fails", func(t *testing.T) {
		router, mockBuscaCEP, mockWeather := setupRouter(t)

		mockBuscaCEP.GetBrasilAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "São Paulo,SP", nil
		}
		mockBuscaCEP.GetViaAPICEPFunc = func(ctx context.Context, cep string) (string, error) {
			return "São Paulo,SP", nil
		}
		mockWeather.GetHGWeatherAPIFunc = func(ctx context.Context, city string) (int, error) {
			return 0, errors.New("weather api error")
		}

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/tempForCep/01001000", nil)
		router.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
		assert.Contains(t, w.Body.String(), "can not find temperature in Celsius for City: São Paulo,SP")
	})
}