package contract

import (
	"encoding/json"
	"fmt"

	"github.com/gzttcydxx/did/models"
	"github.com/hyperledger/fabric-contract-api-go/contractapi"
)

type SmartContract struct {
	contractapi.Contract
}

type ResponseMessage struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

func (s *SmartContract) InitLedger(ctx contractapi.TransactionContextInterface) error {
	args := ctx.GetStub().GetArgs()

	if len(args) != 2 {
		return fmt.Errorf("incorrect number of arguments. Expecting 2")
	}

	chainID := args[1]

	err := ctx.GetStub().PutState("chainID", chainID)
	if err != nil {
		return fmt.Errorf("failed to put chainID to world state: %v", err)
	}

	err = ctx.GetStub().PutState("chainIDverified", []byte("false"))
	if err != nil {
		return fmt.Errorf("failed to put chainIDverified to world state: %v", err)
	}

	return nil
}

func (s *SmartContract) GetChainID(ctx contractapi.TransactionContextInterface) (string, error) {
	chainID, err := ctx.GetStub().GetState("chainID")
	if err != nil {
		return "", fmt.Errorf("failed to get chainID: %v", err)
	}

	return string(chainID), nil
}

func (s *SmartContract) GetChainIDverified(ctx contractapi.TransactionContextInterface) (string, error) {
	chainIDverified, err := ctx.GetStub().GetState("chainIDverified")
	if err != nil {
		return "", fmt.Errorf("failed to get chainIDverified: %v", err)
	}

	return string(chainIDverified), nil
}

func (s *SmartContract) SetChainIDVerified(ctx contractapi.TransactionContextInterface) error {
	err := ctx.GetStub().PutState("chainIDverified", []byte("true"))
	if err != nil {
		return fmt.Errorf("failed to put chainIDverified to world state: %v", err)
	}

	return nil
}

func (s *SmartContract) CreateIdentity(ctx contractapi.TransactionContextInterface, did, publicKey, privateKey string) error {
	readIdentity, _ := s.ReadIdentity(ctx, did)

	if readIdentity != nil {
		return fmt.Errorf("the identity %s already exists", did)
	}

	newdid, err := models.NewDID(did)
	if err != nil {
		return fmt.Errorf("failed to generate DID: %v", err)
	}

	identity, err := models.NewDIDDoc(*newdid, publicKey, privateKey)
	if err != nil {
		return fmt.Errorf("failed to generate DIDdoc: %v", err)
	}

	identityJSON, err := json.Marshal(identity)
	if err != nil {
		return err
	}

	return ctx.GetStub().PutState(did, identityJSON)
}

func (s *SmartContract) ReadIdentity(ctx contractapi.TransactionContextInterface, did string) (*models.DIDDoc, error) {
	DIDdocJSON, err := ctx.GetStub().GetState(did)
	if err != nil {
		return nil, err
	}
	if DIDdocJSON == nil {
		return nil, nil
	}

	var DIDdoc *models.DIDDoc
	err = json.Unmarshal(DIDdocJSON, &DIDdoc)
	if err != nil {
		return nil, err
	}

	return DIDdoc, nil
}

func (s *SmartContract) UpdateIdentity(ctx contractapi.TransactionContextInterface, didDocRaw string) error {
	var newDIDDoc models.DIDDoc

	err := json.Unmarshal([]byte(didDocRaw), &newDIDDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal didDoc: %v", err)
	}

	did := newDIDDoc.ID.Scheme + ":" + newDIDDoc.ID.Method + ":" + newDIDDoc.ID.SpecificID

	readIdentity, err := s.ReadIdentity(ctx, did)
	if err != nil {
		return err
	}
	if readIdentity == nil {
		return fmt.Errorf("the identity %s does not exist", did)
	}

	return ctx.GetStub().PutState(did, []byte(didDocRaw))
}

func (s *SmartContract) DeleteIdentity(ctx contractapi.TransactionContextInterface, did string) error {
	readIdentity, err := s.ReadIdentity(ctx, did)
	if err != nil {
		return err
	}

	if readIdentity == nil {
		return fmt.Errorf("the identity %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}

// 中继链方法
func (s *SmartContract) RegisterCrosschainIdentity(ctx contractapi.TransactionContextInterface, didDocRaw string) error {
	var didDoc models.DIDDoc
	err := json.Unmarshal([]byte(didDocRaw), &didDoc)
	if err != nil {
		return fmt.Errorf("failed to unmarshal didDoc: %v", err)
	}

	chainID := didDoc.ID.ChainID
	if chainID == "" {
		return fmt.Errorf("the chainID is empty")
	}

	result, err := ctx.GetStub().GetState(chainID)
	if err != nil {
		return fmt.Errorf("failed to get chainID: %v", err)
	} else if result != nil {
		return fmt.Errorf("the chainID %s already exists", chainID)
	} else {
		err = ctx.GetStub().PutState(chainID, []byte(didDocRaw))
		if err != nil {
			return fmt.Errorf("failed to put chainID to world state: %v", err)
		}
		return nil
	}
}
