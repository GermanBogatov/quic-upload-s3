package view

type Response struct {
	Code   string      `json:"code"`
	Result interface{} `json:"result"`
	Error  string      `json:"error"`
}
