package models

import (
	"time"

	"github.com/gzttcydxx/did/models"
)

// 订单状态枚举
type OrderStatus int

const (
	Created           OrderStatus = iota // 已创建
	SupplierConfirmed                    // 供应商已确认
	DemanderSelected                     // 买家已选择
	SupplierApproved                     // 供应商已批准
	DemanderApproved                     // 买家已批准
	Completed                            // 已完成
	DemanderCanceled                     // 买家已取消
)

// 订单结构
type Order struct {
	Did            models.DID  `json:"did"`
	DemanderDid    models.DID  `json:"demander_did"`
	SupplierDid    models.DID  `json:"supplier_did"`
	DemandProduct  Part        `json:"demand_product"`
	SupplyProduct  []Product   `json:"supply_product"`
	ComfirmProduct Product     `json:"comfirm_product"`
	Num            int         `json:"num"`
	Status         OrderStatus `json:"status"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}
