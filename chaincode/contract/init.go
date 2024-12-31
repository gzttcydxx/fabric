package contract

import (
	"encoding/json"
	"fmt"

	didModels "github.com/gzttcydxx/did/models"
	"github.com/gzttcydxx/fabric/chaincode/data"
	"github.com/gzttcydxx/fabric/chaincode/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type Orgs []models.Org
type Parts []models.Part
type Products []models.Product

func (s *SmartContract) initOrg(ctx contractapi.TransactionContextInterface) error {
	var orgData Orgs
	err := json.Unmarshal([]byte(data.ORGS), &orgData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal org data: %v", err)
	}

	for _, org := range orgData {
		did, err := didModels.NewDID(fmt.Sprintf("did:org:%s", org.UUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		org.Did = *did

		orgJSON, err := json.Marshal(org)
		if err != nil {
			return fmt.Errorf("failed to marshal org: %v", err)
		}

		// Use ORG_UUID as key to store in blockchain
		err = ctx.GetStub().PutState(org.Did.ToString(), orgJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) initPart(ctx contractapi.TransactionContextInterface) error {
	var partData Parts
	err := json.Unmarshal([]byte(data.PARTS), &partData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal part data: %v", err)
	}

	for _, part := range partData {
		did, err := didModels.NewDID(fmt.Sprintf("did:part:%s", part.UUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		part.Did = *did

		partJSON, err := json.Marshal(part)
		if err != nil {
			return fmt.Errorf("failed to marshal part: %v", err)
		}

		err = ctx.GetStub().PutState(part.Did.ToString(), partJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) initProduct(ctx contractapi.TransactionContextInterface) error {
	var productData Products
	err := json.Unmarshal([]byte(data.PRODUCTS), &productData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal part relation data: %v", err)
	}

	for _, product := range productData {
		did, err := didModels.NewDID(fmt.Sprintf("did:product:%s_%s", product.OrgUUID, product.PartUUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		product.Did = *did

		productJSON, err := json.Marshal(product)
		if err != nil {
			return fmt.Errorf("failed to marshal product: %v", err)
		}

		err = ctx.GetStub().PutState(product.Did.ToString(), productJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}
