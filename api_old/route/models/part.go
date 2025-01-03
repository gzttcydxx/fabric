package models

import (
	"github.com/gzttcydxx/did/models"
)

// Part 零件模型
// @Description 零件信息模型
// @SchemaExample(parts)
type Part struct {
	Did  models.DID `json:"did" swaggertype:"string" example:"did:part:66666666-a0e3-48dc-bff8-7eb7819b7a09"` // 零件DID
	UUID string     `json:"uuid" example:"66666666-a0e3-48dc-bff8-7eb7819b7a09"`                              // 唯一标识符
	Name string     `json:"name" example:"围板"`                                                                // 零件名称
}

// Product 产品模型
// @Description 产品信息模型
// @SchemaExample(products)
type Product struct {
	Did           models.DID `json:"did"`            // 零件关系DID
	OrgUUID       string     `json:"start_id"`       // 供应商ID
	OrgName       string     `json:"start_label"`    // 供应商名称
	PartUUID      string     `json:"end_id"`         // 零件ID
	PartName      string     `json:"end_label"`      // 零件名称
	RelationName  string     `json:"relation_name"`  // 关系名称
	RelationLabel string     `json:"relation_label"` // 关系标签
	Price         float64    `json:"price"`          // 价格
	FitBrand      string     `json:"fit_brand"`      // 适配品牌
	ProductPlace  string     `json:"product_place"`  // 生产地
}
