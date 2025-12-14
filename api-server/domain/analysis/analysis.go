package analysis

import (
	"api-server/domain"
	"context"
	"log"
	"time"
)

type analysisService struct {
	buscaCEPAPIClient domain.BuscaCEPAPIClient
	weatherAPIClient  domain.WeatherAPIClient
	log               *log.Logger
}

func NewAnalysisService(buscaCEPAPIClient domain.BuscaCEPAPIClient, weatherAPIClient domain.WeatherAPIClient,
	log *log.Logger) *analysisService {

	return &analysisService{
		buscaCEPAPIClient: buscaCEPAPIClient,
		weatherAPIClient:  weatherAPIClient,
		log:               log,
	}
}

func (s *analysisService) GetCity(c context.Context, cep string) (string, error) {
	type result struct {
		Source string
		City   string
		Err    error
	}

	// Timeout de 1 segundo para chamada das APIs
	apiCtx, apiCancel := context.WithTimeout(c, 1*time.Second)
	defer apiCancel()

	resultCh := make(chan result, 2)

	// Chamada BrasilAPI
	go func() {
		city, err := s.buscaCEPAPIClient.GetBrasilAPICEP(apiCtx, cep)
		resultCh <- result{Source: "BrasilAPI", City: city, Err: err}
	}()

	// Chamada ViaCEP
	go func() {
		city, err := s.buscaCEPAPIClient.GetViaAPICEP(apiCtx, cep)
		resultCh <- result{Source: "ViaCEP", City: city, Err: err}
	}()

	var lastErr error
	for i := 0; i < 2; i++ {
		select {
		case res := <-resultCh:
			if res.Err == nil {
				s.log.Printf("Resposta recebida da API %s: %s", res.Source, res.City)
				return res.City, nil
			}
			s.log.Printf("Erro ao buscar CEP na API %s: %v", res.Source, res.Err)
			lastErr = res.Err
		case <-apiCtx.Done():
			s.log.Printf("Timeout ao buscar CEP nas APIs")
			return "", apiCtx.Err()
		}
	}

	return "", lastErr
}

func (s *analysisService) GetCelsiusTemperature(c context.Context, city string) (int, error) {
	type result struct {
		Temp int
		Err  error
	}

	// Timeout de 1 segundo para chamada das APIs
	apiCtx, apiCancel := context.WithTimeout(c, 1*time.Second)
	defer apiCancel()

	resultCh := make(chan result, 1)

	// Chamada HG Weather API
	go func() {
		temp, err := s.weatherAPIClient.GetHGWeatherAPI(apiCtx, city)
		resultCh <- result{Temp: temp, Err: err}
	}()

	select {
	case res := <-resultCh:
		if res.Err != nil {
			s.log.Printf("Error Getting Temperature on the API for the City %s: %v", city, res.Err)
			return 0, res.Err
		}
		s.log.Printf("Response from the Temperature API: %v", res.Temp)
		return res.Temp, nil
	case <-apiCtx.Done():
		s.log.Printf("Timeout in API Searching for Temperature")
		return 0, apiCtx.Err()
	}

}
