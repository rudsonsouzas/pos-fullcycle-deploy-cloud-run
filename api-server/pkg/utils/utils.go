package utils

func ConvertCelsiusToFahrenheit(celsius float64) float64 {
	return (celsius * 1.8) + 32
}

func ConvertCelsiusToKelvin(celsius float64) float64 {
	return celsius + 273.15
}

func IsValidCEP(cep string) bool {
	// CEP deve ter 8 dígitos numéricos
	if len(cep) != 8 {
		return false
	}
	for _, c := range cep {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
}
