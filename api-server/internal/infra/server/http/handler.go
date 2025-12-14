package http

import (
	"log"
	"os"

	"api-server/domain"

	"github.com/gin-gonic/gin"
)

type handler struct {
	analisysService domain.AnalysisService
	log             *log.Logger
}

func NewHandler(analisysService domain.AnalysisService, log *log.Logger,
) *gin.Engine {
	handler := &handler{
		analisysService: analisysService,
		log:             log,
	}

	gin.SetMode(gin.ReleaseMode)

	router := gin.New()
	router.Use(gin.Recovery())

	if os.Getenv("ENV") == "local" {
		gin.SetMode(gin.DebugMode)
		router.Use(gin.Logger())
	}

	router.GET("/tempForCep/:cep", handler.RunAnalysis)

	return router
}
