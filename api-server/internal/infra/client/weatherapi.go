package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"time"

	"api-server/domain"
	httpclient "api-server/pkg/http_client"

	"github.com/cenkalti/backoff"
)

type WeatherAPIClient struct {
	httpClient httpclient.HTTPClient
	log        *log.Logger
	apiKey     string
}

func NewWeatherAPIClient(httpClient httpclient.HTTPClient, log *log.Logger, apiKey string) *WeatherAPIClient {
	return &WeatherAPIClient{
		httpClient: httpClient,
		log:        log,
		apiKey:     apiKey,
	}
}

func (awc *WeatherAPIClient) getTemperature(ctx context.Context, url string) ([]byte, error) {
	var resBody []byte

	ebo := backoff.NewExponentialBackOff()
	ebo.MaxInterval = 1 * time.Second

	if err := backoff.Retry(func() error {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			awc.log.Printf("error to create request to search Temperature through URL: %s. [Error]: %s", url, err.Error())
			return err
		}

		req.Header.Add("Content-Type", "application/json")

		res, err := awc.httpClient.Do(req)
		if err != nil {
			awc.log.Printf("error to search Temperature through URL: %s. [Error]: %s", url, err.Error())
			return err
		}
		defer func() {
			err = res.Body.Close()
			if err != nil {
				awc.log.Printf("error to close response body from URL: %s. [Error]: %s", url, err.Error())
				return
			}
		}()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			awc.log.Printf("error to read response body from URL: %s. [Error]: %s", url, err.Error())
			return err
		}

		if res.StatusCode != 200 {
			awc.log.Printf("Search Temperature through URL [%s]- API status code %d: %s", url, res.StatusCode, bodyBytes)
			return err
		}

		resBody = bodyBytes
		return nil

	}, backoff.WithContext(backoff.WithMaxRetries(ebo, uint64(5)), ctx)); err != nil {
		awc.log.Printf("error to search Temperature through URL: %s. [Error]: %s", url, err.Error())
		return []byte{}, err
	}

	return resBody, nil
}

func (awc *WeatherAPIClient) GetHGWeatherAPI(ctx context.Context, city string) (int, error) {
	var weatherAPIResponse *domain.HGWeatherAPIResponse

	baseURL := "https://api.hgbrasil.com/weather"
	params := url.Values{}
	params.Add("city_name", city)
	params.Add("key", awc.apiKey)
	weatherAPIUrl := baseURL + "?" + params.Encode()

	resBody, err := awc.getTemperature(ctx, weatherAPIUrl)
	if err != nil {
		awc.log.Printf("error on get info from HG WeatherAPI to the city: %s. [Erro]: %s", city, err.Error())
		return 0, err
	}

	if err := json.Unmarshal(resBody, &weatherAPIResponse); err != nil {
		awc.log.Printf("error to translate response Body to HG WeatherAPI pattern: %s", err.Error())
		return 0, err
	}

	if !weatherAPIResponse.ValidKey {
		awc.log.Printf("invalid API Key provided for HG WeatherAPI")
		return 0, fmt.Errorf("invalid API Key provided for HG WeatherAPI")
	}

	return weatherAPIResponse.Results.Temp, nil
}
