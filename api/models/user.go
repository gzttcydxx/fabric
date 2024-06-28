package models

import (
	"github.com/gzttcydxx/did/models"
)

type User struct {
	Did  models.DID `json:"did"`
	Name string     `json:"name"`
	Role string     `json:"role"`
	Org  Org        `json:"org"`
}
