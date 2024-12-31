package sdk

import (
	"github.com/gzttcydxx/api/utils/crud"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-gateway/pkg/client"
)

// 组织服务
type OrgService struct {
	*crud.CRUDService[*models.Org]
}

func NewOrgService(contract *client.Contract) *OrgService {
	return &OrgService{
		CRUDService: crud.NewCRUDService[*models.Org](
			contract,
			"org",
			crud.CRUDMethods{
				Create: "CreateOrg",
				Query:  "QueryOrgs",
				Read:   "ReadOrg",
				Update: "UpdateOrg",
				Delete: "DeleteOrg",
			},
		),
	}
}

// 零件服务
type PartService struct {
	*crud.CRUDService[*models.Part]
}

func NewPartService(contract *client.Contract) *PartService {
	return &PartService{
		CRUDService: crud.NewCRUDService[*models.Part](
			contract,
			"part",
			crud.CRUDMethods{
				Create: "CreatePart",
				Query:  "QueryParts",
				Read:   "ReadPart",
				Update: "UpdatePart",
				Delete: "DeletePart",
			},
		),
	}
}

// 产品服务
type ProductService struct {
	*crud.CRUDService[*models.Product]
}

func NewProductService(contract *client.Contract) *ProductService {
	return &ProductService{
		CRUDService: crud.NewCRUDService[*models.Product](
			contract,
			"product",
			crud.CRUDMethods{
				Create: "CreateProduct",
				Query:  "QueryProducts",
				Read:   "ReadProduct",
				Update: "UpdateProduct",
				Delete: "DeleteProduct",
			},
		),
	}
}
