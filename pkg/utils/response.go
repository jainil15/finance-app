package utils

type Response struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}
