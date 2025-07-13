package api

import (
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

// validCurrency validates if the currency is one of the supported currencies
var validCurrency = validator.Func(func(fieldLevel validator.FieldLevel) bool {
	if currency, ok := fieldLevel.Field().Interface().(string); ok {
		return isSupportedCurrency(currency)
	}
	return false
})

// isSupportedCurrency checks if the given currency is supported
func isSupportedCurrency(currency string) bool {
	switch currency {
	case "USD", "EUR", "UAH":
		return true
	}
	return false
}

// RegisterValidators registers custom validators with gin
func RegisterValidators() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
}