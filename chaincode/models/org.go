package models

import (
	"github.com/gzttcydxx/did/models"
)

type Org struct {
	Did  models.DID `json:"did"`  // 组织DID
	Name string     `json:"name"` // 组织名称
}
