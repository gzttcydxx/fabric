package models

import "github.com/gzttcydxx/did/models"

// Part 零件模型
type Part struct {
	Did  models.DID `json:"did" required:"false"`
	UUID string     `json:"uuid" required:"false" example:"0006321f-a0e3-48dc-bff8-7eb7819b7a09"`
	Name string     `json:"name" required:"false" example:"围板"`
}
