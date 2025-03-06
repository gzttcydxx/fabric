package models

import (
	"time"

	"github.com/gzttcydxx/did/models"
)

// 订单状态枚举
type OrderStatus int

const (
	Created         OrderStatus = iota // 已创建(买方已选择卖家)
	SellerConfirmed                    // 卖家已确认
	BuyerConfirmed                     // 买家已确认
	SellerCanceled                     // 卖家已取消
	BuyerCanceled                      // 买家已取消
	Completed                          // 已完成
)

// 订单结构
type Order struct {
	Did       models.DID  `json:"did" required:"false"`
	BuyerDid  models.DID  `json:"buyer_did" required:"false"`
	SellerDid models.DID  `json:"seller_did" required:"false"`
	Product   Product     `json:"product" required:"false"`
	Num       int         `json:"num" required:"false"`
	Status    OrderStatus `json:"status" required:"false" example:"0" doc:"0:已创建,1:供应商已确认,2:买家已选择,3:供应商已批准,4:买家已批准,5:供应商已取消,6:买家已取消,7:已完成"`
	CreatedAt time.Time   `json:"created_at" required:"false" example:"2024-01-01T00:00:00Z"`
	UpdatedAt time.Time   `json:"updated_at" required:"false" example:"2024-01-01T00:00:00Z"`
}

type OrderResponse struct {
	Status   Status `json:"status"`
	OrderDid string `json:"order_did"`
}
