package models

import (
	"github.com/gzttcydxx/did/models"
)

type User struct {
	Did  models.DID `json:"did"`  // 用户DID
	Name string     `json:"name"` // 用户名称
	Role string     `json:"role"` // 用户角色
	Org  Org        `json:"org"`  // 用户所属组织
}
