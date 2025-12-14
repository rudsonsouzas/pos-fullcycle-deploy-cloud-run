package main

import (
	"os"
	"time"

	"os/signal"
	"syscall"

	"api-server/domain/analysis"
	"api-server/internal/infra/client"
	"api-server/internal/infra/server/http"
	"api-server/pkg/env"
	httpclient "api-server/pkg/http_client"
	"log"
)

const (
	envApplicationPort = "APP_PORT"
	envWeatherAPIKey   = "WEATHER_API_KEY"

	defaultApplicationPort = "8080"
)

var (
	version, build, date string
)

func main() {
	logger := log.New(os.Stdout, "api-server-temperature - ", log.LstdFlags)

	env.CheckRequired(logger, envWeatherAPIKey)

	logger.Printf("API Busca Temperatura com CEP - version:%s; build:%s; date:%s", version, build, date)

	buscaCEPAPIClient := client.NewBuscaCEPAPIClient(httpclient.NewHTTPClient(60*time.Second), logger)

	weatherAPIClient := client.NewWeatherAPIClient(httpclient.NewHTTPClient(60*time.Second), logger, getWeatherAPIKey())

	analysisService := analysis.NewAnalysisService(buscaCEPAPIClient, weatherAPIClient, logger)

	handler := http.NewHandler(analysisService, logger)

	/*
	 * Server...
	 */
	server := http.New(getApplicationPort(), handler, logger)
	server.ListenAndServe()

	/*
	 * Graceful shutdown...
	 */
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan
	server.Shutdown()
}

func getApplicationPort() string {
	return env.GetString(envApplicationPort, defaultApplicationPort)
}

func getWeatherAPIKey() string {
	return env.GetString(envWeatherAPIKey)
}
