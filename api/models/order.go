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
	SupplierCanceled                     // 供应商已取消
	DemanderCanceled                     // 买家已取消
	Completed                            // 已完成
)

// 订单结构
type Order struct {
	Did            models.DID  `json:"did" required:"false"`
	DemanderDid    models.DID  `json:"demander_did" required:"false"`
	SupplierDid    models.DID  `json:"supplier_did" required:"false"`
	DemandProduct  Part        `json:"demand_product" required:"false"`
	SupplyProduct  []Product   `json:"supply_product" required:"false"`
	ComfirmProduct Product     `json:"comfirm_product" required:"false"`
	Num            int         `json:"num" required:"false"`
	Status         OrderStatus `json:"status" required:"false" example:"0" doc:"0:已创建,1:供应商已确认,2:买家已选择,3:供应商已批准,4:买家已批准,5:供应商已取消,6:买家已取消,7:已完成"`
	CreatedAt      time.Time   `json:"created_at" required:"false" example:"2024-01-01T00:00:00Z"`
	UpdatedAt      time.Time   `json:"updated_at" required:"false" example:"2024-01-01T00:00:00Z"`
}

type OrderProductInput struct {
	GetInput
	Body Product `json:"body" required:"false"`
}
