package models

type QRResponse struct {
	Q          [][]float64 `json:"q"`
	R          [][]float64 `json:"r"`
	Statistics Statistics  `json:"statistics"`
}

type APIError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error APIError `json:"error"`
}
