package domain

import "context"

type AnalysisService interface {
	GetCity(c context.Context, cep string) (string, error)
	GetCelsiusTemperature(c context.Context, city string) (int, error)
}
