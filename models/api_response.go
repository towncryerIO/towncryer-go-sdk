package models

type ApiResponse[T any] struct {
	Code    string   `json:"code"`
	Message string   `json:"message"`
	Errors  []string `json:"errors"`
	Data    T        `json:"data"`
}
