package analysis

import (
	"context"
	"errors"
	"log"
	"os"
	"testing"
	"time"

	"api-server/domain/mocks"
	"github.com/stretchr/testify/assert"
)

func TestAnalysisService_GetCity(t *testing.T) {
	logger := log.New(os.Stdout, "test - ", log.LstdFlags)

	t.Run("should return city from BrasilAPI first", func(t *testing.T) {
		mockBuscaCEPClient := &mocks.MockBuscaCEPAPIClient{
			GetBrasilAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				return "City From BrasilAPI", nil
			},
			GetViaAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				// This should not be called
				return "", errors.New("should not be called")
			},
		}
		service := NewAnalysisService(mockBuscaCEPClient, nil, logger)

		city, err := service.GetCity(context.Background(), "12345678")

		assert.NoError(t, err)
		assert.Equal(t, "City From BrasilAPI", city)
	})

	t.Run("should return city from ViaCEP when BrasilAPI fails", func(t *testing.T) {
		mockBuscaCEPClient := &mocks.MockBuscaCEPAPIClient{
			GetBrasilAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				return "", errors.New("brasil api error")
			},
			GetViaAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				return "City From ViaCEP", nil
			},
		}
		service := NewAnalysisService(mockBuscaCEPClient, nil, logger)

		city, err := service.GetCity(context.Background(), "12345678")

		assert.NoError(t, err)
		assert.Equal(t, "City From ViaCEP", city)
	})

	t.Run("should timeout when both APIs are slow", func(t *testing.T) {
		mockBuscaCEPClient := &mocks.MockBuscaCEPAPIClient{
			GetBrasilAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				time.Sleep(2 * time.Second)
				return "City From BrasilAPI", nil
			},
			GetViaAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				time.Sleep(2 * time.Second)
				return "City From ViaCEP", nil
			},
		}
		service := NewAnalysisService(mockBuscaCEPClient, nil, logger)

		_, err := service.GetCity(context.Background(), "12345678")

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("should return error when both APIs fail", func(t *testing.T) {
		mockBuscaCEPClient := &mocks.MockBuscaCEPAPIClient{
			GetBrasilAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				return "", errors.New("brasil api error")
			},
			GetViaAPICEPFunc: func(ctx context.Context, cep string) (string, error) {
				return "", errors.New("via cep error")
			},
		}
		service := NewAnalysisService(mockBuscaCEPClient, nil, logger)

		_, err := service.GetCity(context.Background(), "12345678")

		assert.Error(t, err)
	})
}

func TestAnalysisService_GetCelsiusTemperature(t *testing.T) {
	logger := log.New(os.Stdout, "test - ", log.LstdFlags)

	t.Run("should return temperature successfully", func(t *testing.T) {
		mockWeatherClient := &mocks.MockWeatherAPIClient{
			GetHGWeatherAPIFunc: func(ctx context.Context, city string) (int, error) {
				return 25, nil
			},
		}
		service := NewAnalysisService(nil, mockWeatherClient, logger)

		temp, err := service.GetCelsiusTemperature(context.Background(), "São Paulo")

		assert.NoError(t, err)
		assert.Equal(t, 25, temp)
	})

	t.Run("should timeout when api is slow", func(t *testing.T) {
		mockWeatherClient := &mocks.MockWeatherAPIClient{
			GetHGWeatherAPIFunc: func(ctx context.Context, city string) (int, error) {
				time.Sleep(2 * time.Second)
				return 25, nil
			},
		}
		service := NewAnalysisService(nil, mockWeatherClient, logger)

		_, err := service.GetCelsiusTemperature(context.Background(), "São Paulo")

		assert.Error(t, err)
		assert.Equal(t, context.DeadlineExceeded, err)
	})

	t.Run("should return error when api fails", func(t *testing.T) {
		mockWeatherClient := &mocks.MockWeatherAPIClient{
			GetHGWeatherAPIFunc: func(ctx context.Context, city string) (int, error) {
				return 0, errors.New("weather api error")
			},
		}
		service := NewAnalysisService(nil, mockWeatherClient, logger)

		_, err := service.GetCelsiusTemperature(context.Background(), "São Paulo")

		assert.Error(t, err)
	})
}
