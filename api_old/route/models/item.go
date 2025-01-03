package models

import (
	"github.com/gzttcydxx/did/models"
)

type Item struct {
	Did   models.DID `json:"did"`   // 物品ID
	Name  string     `json:"name"`  // 物品名称
	Type  ItemType   `json:"type"`  // 物品类型
	Org   Org        `json:"org"`   // 物品所属组织
	Owner User       `json:"owner"` // 物品所有者
	Price int        `json:"price"` // 物品价格
}

type ItemType struct {
	Did  models.DID `json:"did"`  // 物品类型ID
	Name string     `json:"name"` // 物品类型名称
	Unit string     `json:"unit"` // 物品单位
}

type ItemDemand struct {
	ItemType
	Num int `json:"num"` // 需求数量
}

type ItemStock struct {
	Item
	Num int `json:"num"` // 库存数量
}
