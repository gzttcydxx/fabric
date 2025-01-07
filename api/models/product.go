package models

import "github.com/gzttcydxx/did/models"

// Product 产品模型
type Product struct {
	Did           models.DID `json:"did" required:"false"`
	OrgUUID       string     `json:"start_id" required:"false"`
	OrgName       string     `json:"start_label" required:"false"`
	PartUUID      string     `json:"end_id" required:"false"`
	PartName      string     `json:"end_label" required:"false"`
	RelationName  string     `json:"relation_name" required:"false"`
	RelationLabel string     `json:"relation_label" required:"false"`
	Price         float64    `json:"price" required:"false"`
	FitBrand      string     `json:"fit_brand" required:"false"`
	ProductPlace  string     `json:"product_place" required:"false"`
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
