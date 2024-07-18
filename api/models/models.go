package models

import (
	did "github.com/gzttcydxx/did/models"
)

type Org struct {
	Did  did.DID `json:"did"`
	Name string  `json:"name"`
}

type User struct {
	Did  did.DID `json:"did"`
	Name string  `json:"name"`
	Role string  `json:"role"`
	Org  Org     `json:"org"`
}
