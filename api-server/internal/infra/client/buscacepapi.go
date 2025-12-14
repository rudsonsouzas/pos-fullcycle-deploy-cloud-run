package client

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"api-server/domain"
	httpclient "api-server/pkg/http_client"

	"github.com/cenkalti/backoff"
)

type BuscaCEPAPIClient struct {
	httpClient httpclient.HTTPClient
	log        *log.Logger
}

// type BuscaCEPAPI interface {
// 	GetDolarQuote(c context.Context) (*domain.BuscaCEPAPIResponse, error)
// }

func NewBuscaCEPAPIClient(httpClient httpclient.HTTPClient, log *log.Logger) *BuscaCEPAPIClient {
	return &BuscaCEPAPIClient{
		httpClient: httpClient,
		log:        log,
	}
}

func (awc *BuscaCEPAPIClient) getCEP(ctx context.Context, url string) ([]byte, error) {
	var resBody []byte

	ebo := backoff.NewExponentialBackOff()
	ebo.MaxInterval = 1 * time.Second

	if err := backoff.Retry(func() error {

		req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
		if err != nil {
			awc.log.Printf("error to create request to search CEP through URL: %s. [Error]: %s", url, err.Error())
			return err
		}

		req.Header.Add("Content-Type", "application/json")

		res, err := awc.httpClient.Do(req)
		if err != nil {
			awc.log.Printf("error to search CEP through URL: %s. [Error]: %s", url, err.Error())
			return err
		}
		defer func() {
			err = res.Body.Close()
			if err != nil {
				awc.log.Printf("error to close response body from URL: %s. [Erro]: %s", url, err.Error())
				return
			}
		}()

		bodyBytes, err := io.ReadAll(res.Body)
		if err != nil {
			awc.log.Printf("error to read response body from URL: %s. [Error]: %s", url, err.Error())
			return err
		}

		if res.StatusCode != 200 {
			awc.log.Printf("Search CEP through URL [%s]- API status code %d: %s", url, res.StatusCode, bodyBytes)
			return err
		}

		resBody = bodyBytes
		return nil

	}, backoff.WithContext(backoff.WithMaxRetries(ebo, uint64(5)), ctx)); err != nil {
		awc.log.Printf("error to search CEP through URL: %s. [Error]: %s", url, err.Error())
		return []byte{}, err
	}

	return resBody, nil
}

func (awc *BuscaCEPAPIClient) GetBrasilAPICEP(ctx context.Context, cep string) (string, error) {
	var brasilAPIResponse *domain.BrasilAPIResponse
	brasilAPIUrl := "https://brasilapi.com.br/api/cep/v1/" + cep

	resBody, err := awc.getCEP(ctx, brasilAPIUrl)
	if err != nil {
		awc.log.Printf("error to search CEP in BrasilAPI: %s. [Erro]: %s", cep, err.Error())
		return "", err
	}

	if err := json.Unmarshal(resBody, &brasilAPIResponse); err != nil {
		awc.log.Printf("error to translate response Body to BrasilAPI pattern: %s", err.Error())
		return "", err
	}

	cityInfo := fmt.Sprint(brasilAPIResponse.City, ",", brasilAPIResponse.State)

	return cityInfo, nil
}

func (awc *BuscaCEPAPIClient) GetViaAPICEP(ctx context.Context, cep string) (string, error) {
	var viaAPIResponse *domain.ViaCEPAPIResponse
	viaAPIUrl := "https://viacep.com.br/ws/" + cep + "/json/"

	resBody, err := awc.getCEP(ctx, viaAPIUrl)
	if err != nil {
		awc.log.Printf("error to search CEP in ViaAPI: %s. [Erro]: %s", cep, err.Error())
		return "", err
	}

	if err := json.Unmarshal(resBody, &viaAPIResponse); err != nil {
		awc.log.Printf("error to translate response Body to ViaAPI pattern: %s", err.Error())
		return "", err
	}

	if viaAPIResponse.Erro != "" {
		awc.log.Printf("error returned by ViaAPI for CEP: %s. [Erro]: %s", cep, viaAPIResponse.Erro)
		return "", fmt.Errorf("unable to find CEP in ViaAPI")
	}

	cityInfo := fmt.Sprint(viaAPIResponse.Localidade, ",", viaAPIResponse.Uf)

	return cityInfo, nil
}
