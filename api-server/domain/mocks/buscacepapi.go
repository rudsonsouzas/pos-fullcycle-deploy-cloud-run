package mocks

import (
	"context"
	"time"
)

type MockBuscaCEPAPIClient struct {
	GetBrasilAPICEPFunc func(ctx context.Context, cep string) (string, error)
	GetViaAPICEPFunc    func(ctx context.Context, cep string) (string, error)
}

func (m *MockBuscaCEPAPIClient) GetBrasilAPICEP(ctx context.Context, cep string) (string, error) {
	// Simulate delay if needed for timeout tests
	if delay := ctx.Value("delay"); delay != nil {
		time.Sleep(delay.(time.Duration))
	}
	return m.GetBrasilAPICEPFunc(ctx, cep)
}

func (m *MockBuscaCEPAPIClient) GetViaAPICEP(ctx context.Context, cep string) (string, error) {
	// Simulate delay if needed for timeout tests
	if delay := ctx.Value("delay"); delay != nil {
		time.Sleep(delay.(time.Duration))
	}
	return m.GetViaAPICEPFunc(ctx, cep)
}
