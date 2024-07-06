package encryptinterceptor

type EncryptData struct {
	data string
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
