package models

import (
	"github.com/gzttcydxx/did/models"
)

type Org struct {
	Did  models.DID `json:"did"`
	Name string     `json:"name"`
}
