package utils

import "strings"

// List of supported currencies
var SupportedCurrencies = []string{
	"USD",
	"NGN",
	"GHS"}

// Check if a currency is valid
func IsSupportedCurrency(currency string) bool {
	for _, c := range SupportedCurrencies {
		if strings.ToUpper(currency) == c {
			return true
		}
	}
	return false
}
