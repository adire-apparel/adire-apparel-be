package constants

type Currency struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

var Currencies = []Currency{
	{Code: "USD", Name: "US Dollar"},
	{Code: "EUR", Name: "Euro"},
	{Code: "GBP", Name: "British Pound"},
	{Code: "NGN", Name: "Nigerian Naira"},
	{Code: "CAD", Name: "Canadian Dollar"},
	{Code: "AUD", Name: "Australian Dollar"},
	{Code: "CHF", Name: "Swiss Franc"},
	{Code: "CNY", Name: "Chinese Yuan"},
}

func GetCurrencyByCode(code string) (*Currency, bool) {
	for _, currency := range Currencies {
		if currency.Code == code {
			return &currency, true
		}
	}
	return nil, false
}

func GetCurrencyCodes() []string {
	codes := make([]string, len(Currencies))
	for i, currency := range Currencies {
		codes[i] = currency.Code
	}
	return codes
}

func IsValidCurrency(code string) bool {
	_, found := GetCurrencyByCode(code)
	return found
}
