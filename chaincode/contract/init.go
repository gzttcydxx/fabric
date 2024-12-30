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
type PartRelations []models.PartRelation

func (s *SmartContract) initOrg(ctx contractapi.TransactionContextInterface) error {
	var orgData Orgs
	err := json.Unmarshal([]byte(data.ORGS), &orgData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal org data: %v", err)
	}

	for _, org := range orgData {
		did, err := didModels.NewDID(fmt.Sprintf("did:uuid:%s", org.UUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		org.Did = *did

		orgJSON, err := json.Marshal(org)
		if err != nil {
			return fmt.Errorf("failed to marshal org: %v", err)
		}

		// Use ORG_UUID as key to store in blockchain
		err = ctx.GetStub().PutState(fmt.Sprintf("ORG_%s", org.Did.ToString()), orgJSON)
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
		did, err := didModels.NewDID(fmt.Sprintf("did:uuid:%s", part.UUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		part.Did = *did

		partJSON, err := json.Marshal(part)
		if err != nil {
			return fmt.Errorf("failed to marshal part: %v", err)
		}

		err = ctx.GetStub().PutState(fmt.Sprintf("PART_%s", part.Did.ToString()), partJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}

func (s *SmartContract) initPartRelation(ctx contractapi.TransactionContextInterface) error {
	var partRelationData PartRelations
	err := json.Unmarshal([]byte(data.PRODUCTS), &partRelationData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal part relation data: %v", err)
	}

	for _, partRelation := range partRelationData {
		did, err := didModels.NewDID(fmt.Sprintf("did:uuid:%s_%s", partRelation.OrgUUID, partRelation.PartUUID))
		if err != nil {
			return fmt.Errorf("failed to create DID: %v", err)
		}
		partRelation.Did = *did

		partRelationJSON, err := json.Marshal(partRelation)
		if err != nil {
			return fmt.Errorf("failed to marshal part relation: %v", err)
		}

		err = ctx.GetStub().PutState(fmt.Sprintf("PRODUCT_%s", partRelation.Did.ToString()), partRelationJSON)
		if err != nil {
			return fmt.Errorf("failed to put state: %v", err)
		}
	}

	return nil
}
