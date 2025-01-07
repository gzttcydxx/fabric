package models

import (
	"github.com/gzttcydxx/did/models"
)

type Part struct {
	Did  models.DID `json:"did"`  // 零件DID
	UUID string     `json:"uuid"` // 唯一标识符
	Name string     `json:"name"` // 零件名称
}

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

func (p *Product) OrgDid() models.DID {
	did, err := models.NewDID("did:org:" + p.OrgUUID)
	if err != nil {
		return models.DID{}
	}
	return *did
}

func (p *Product) PartDid() models.DID {
	did, err := models.NewDID("did:part:" + p.PartUUID)
	if err != nil {
		return models.DID{}
	}
	return *did
}
