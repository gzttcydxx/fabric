package models

import (
	"time"

	"github.com/gzttcydxx/did/models"
)

type TransactionType string

const (
	START  TransactionType = "START"  // 交易开始
	PADING TransactionType = "PADING" // 交易等待
	ACCEPT TransactionType = "ACCEPT" // 交易接受
	REJECT TransactionType = "REJECT" // 交易拒绝
	END    TransactionType = "END"    // 交易结束
)

type Transaction struct {
	Did    models.DID      `json:"did"`    // 交易ID
	Demand models.DID      `json:"demand"` // 需求方ID
	Supply models.DID      `json:"supply"` // 供应方ID
	Type   string          `json:"type"`   // 交易类型
	Amount int             `json:"amount"` // 交易数量
	Item   Item            `json:"item"`   // 交易物品
	Status TransactionType `json:"status"` // 交易状态
	Time   time.Time       `json:"time"`   // 交易时间
	// StartTime  time.Time       `json:"starttime"`  // 发起交易时间
	// AcceptTime time.Time       `json:"accepttime"` // 接受交易时间
	// EndTime    time.Time       `json:"endtime"`    // 完成交易时间
}
