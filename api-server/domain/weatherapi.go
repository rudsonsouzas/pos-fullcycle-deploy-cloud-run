package domain

import "context"

type HGWeatherAPIResponse struct {
	ValidKey bool             `json:"valid_key"`
	Results  HGWeatherResults `json:"results"`
}

type HGWeatherResults struct {
	Temp int    `json:"temp"`
	City string `json:"city"`
}

type WeatherAPIClient interface {
	GetHGWeatherAPI(ctx context.Context, city string) (int, error)
}
