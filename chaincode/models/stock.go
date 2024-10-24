package models

import (
	"fmt"

	"github.com/gzttcydxx/did/models"
)

type Stock struct {
	Did   models.DID            `json:"did"`   // 库存ID
	Items map[string]*ItemStock `json:"items"` // 物品ID到数量的映射
}

// 创建一个新的库存
func NewStock(did models.DID) *Stock {
	return &Stock{
		Did:   did,
		Items: make(map[string]*ItemStock),
	}
}

// 添加物品到库存
func (s *Stock) AddItem(item Item, num int) error {
	if num <= 0 {
		return fmt.Errorf("数量必须大于 0")
	}
	itemDid := item.Did.ToString()
	if _, exists := s.Items[itemDid]; exists {
		return fmt.Errorf("物品 %s 已经存在于库存中", itemDid)
	}
	s.Items[itemDid] = &ItemStock{
		Item: &item,
		Num:  num,
	}
	return nil
}

// 更新库存中物品的数量
func (s *Stock) UpdateItem(item Item, changeNum int) error {
	itemDid := item.Did.ToString()
	itemStock, exists := s.Items[itemDid]
	if !exists {
		return fmt.Errorf("物品不存在于库存中")
	}
	num := itemStock.Num
	num += changeNum
	if num < 0 {
		return fmt.Errorf("库存数量不足")
	} else if num == 0 {
		delete(s.Items, itemDid)
	} else {
		itemStock.Num = num
		s.Items[itemDid] = itemStock
	}
	return nil
}
