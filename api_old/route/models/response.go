package models

type Response struct {
	StatusCode int         `json:"statusCode" example:"200"`
	Message    interface{} `json:"message"`
}
