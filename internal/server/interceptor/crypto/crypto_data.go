package cryptointerceptor

type CryptoData struct {
	Data string `json:"data"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
