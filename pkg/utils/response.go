package utils

type Response struct {
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

type Error struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error"`
}
