package mocks

import (
	"context"
	"time"
)

type MockWeatherAPIClient struct {
	GetHGWeatherAPIFunc func(ctx context.Context, city string) (int, error)
}

func (m *MockWeatherAPIClient) GetHGWeatherAPI(ctx context.Context, city string) (int, error) {
	// Simulate delay if needed for timeout tests
	if delay := ctx.Value("delay"); delay != nil {
		time.Sleep(delay.(time.Duration))
	}
	return m.GetHGWeatherAPIFunc(ctx, city)
}
