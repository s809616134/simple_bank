package util

const (
	USD = "USD"
	UER = "EUR"
	CAD = "CAD"
)

func IsSupportedCurrency(currency string) bool {
	switch currency {
	case USD, UER, CAD:
		return true
	}
	return false
}
