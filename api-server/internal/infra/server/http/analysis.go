package http

import (
	"api-server/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) RunAnalysis(c *gin.Context) {
	cep := c.Param("cep")
	if !utils.IsValidCEP(cep) {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": "invalid zipcode"})
		return
	}

	// Call the analysis service to run
	cityInfo, err := h.analisysService.GetCity(c.Request.Context(), cep)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "can not find zipcode"})
		return
	}

	celsiusTemp, err := h.analisysService.GetCelsiusTemperature(c.Request.Context(), cityInfo)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			gin.H{"error": "can not find temperature in Celsius for City: " + cityInfo + ".", "details": err.Error()})
		return
	}

	celsiusTempFloat := float64(celsiusTemp)
	fahrenheitTemp := utils.ConvertCelsiusToFahrenheit(celsiusTempFloat)
	kelvinTemp := utils.ConvertCelsiusToKelvin(celsiusTempFloat)

	// Arredondar para 2 casas decimais para consistência no JSON
	fahrenheitTemp = float64(int(fahrenheitTemp*100)) / 100
	kelvinTemp = float64(int(kelvinTemp*100)) / 100

	c.JSON(http.StatusOK, gin.H{"temp_C": celsiusTemp, "temp_F": fahrenheitTemp, "temp_K": kelvinTemp})
	c.Next()
}
