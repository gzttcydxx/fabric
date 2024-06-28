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

//////////////////////////////////////////////////////////////

func (s *SmartContract) ReadUser(ctx contractapi.TransactionContextInterface, did string) (User, error) {
	var user User
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return user, fmt.Errorf("failed to get user info: %v", err)
	}
	if result == nil {
		return user, fmt.Errorf("the user %s does not exist", did)
	}

	err = json.Unmarshal(result, &user)
	if err != nil {
		return user, fmt.Errorf("failed to unmarshal user info: %v", err)
	}

	return user, nil
}

func (s *SmartContract) ReadUsers(ctx contractapi.TransactionContextInterface) ([]User, error) {
	queryString := `{"selector":{"role":{"$exists":true}}}`

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var users []User
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var user User
		err = json.Unmarshal(queryResponse.Value, &user)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal user: %v", err)
		}
		users = append(users, user)
	}

	return users, nil
}

func (s *SmartContract) CreateUser(ctx contractapi.TransactionContextInterface, userJson string) error {
	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	did := user.Did
	org := user.Org

	userResult, _ := s.ReadUser(ctx, did.ToString())
	if userResult.Did.SpecificID != "" {
		return fmt.Errorf("the user %s already exists", did)
	}
	// if err != nil {
	// 	return fmt.Errorf("failed to get user info: %v", err)
	// }

	orgResult, err := s.ReadOrg(ctx, org.Did.ToString())
	if orgResult.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", org.Did)
	}
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(userJson))
	if err != nil {
		return fmt.Errorf("failed to put user: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateUser(ctx contractapi.TransactionContextInterface, userJson string) error {
	var user User
	err := json.Unmarshal([]byte(userJson), &user)
	if err != nil {
		return fmt.Errorf("failed to unmarshal user: %v", err)
	}

	did := user.Did

	readUser, err := s.ReadUser(ctx, did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if readUser.Did.SpecificID == "" {
		return fmt.Errorf("the user %s does not exist", did)
	}

	if user.Name != "" {
		readUser.Name = user.Name
	}
	if user.Role != "" {
		readUser.Role = user.Role
	}
	if user.Org.Did.SpecificID != "" {
		orgResult, err := s.ReadOrg(ctx, user.Org.Did.ToString())
		if err != nil {
			return fmt.Errorf("failed to get org info: %v", err)
		}
		if orgResult.Did.SpecificID == "" {
			return fmt.Errorf("the org %s does not exist", user.Org.Did)
		}
		readUser.Org = user.Org
	}

	readJson, err := json.Marshal(readUser)
	if err != nil {
		return fmt.Errorf("failed to marshal user: %v", err)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(readJson))
	if err != nil {
		return fmt.Errorf("failed to put user: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteUser(ctx contractapi.TransactionContextInterface, did string) error {
	readUser, err := s.ReadUser(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get user info: %v", err)
	}
	if readUser.Did.SpecificID == "" {
		return fmt.Errorf("the user %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}

/////////////////////////////////////////////////////////////

func (s *SmartContract) ReadOrg(ctx contractapi.TransactionContextInterface, did string) (Org, error) {
	var org Org
	result, err := ctx.GetStub().GetState(did)
	if err != nil {
		return org, fmt.Errorf("failed to get org info: %v", err)
	}
	if result == nil {
		return org, fmt.Errorf("the org %s does not exist", did)
	}

	err = json.Unmarshal(result, &org)
	if err != nil {
		return org, fmt.Errorf("failed to unmarshal org info: %v", err)
	}

	return org, nil
}

func (s *SmartContract) ReadOrgs(ctx contractapi.TransactionContextInterface) ([]Org, error) {
	queryString := `{"selector":{"name":{"$exists":true},"role":{"$exists":false}}}`

	iterator, err := ctx.GetStub().GetQueryResult(queryString)
	if err != nil {
		return nil, fmt.Errorf("failed to get query result: %v", err)
	}
	defer iterator.Close()

	var orgs []Org
	for iterator.HasNext() {
		queryResponse, err := iterator.Next()
		if err != nil {
			return nil, fmt.Errorf("failed to get next query result: %v", err)
		}

		var org Org
		err = json.Unmarshal(queryResponse.Value, &org)
		if err != nil {
			return nil, fmt.Errorf("failed to unmarshal org: %v", err)
		}
		orgs = append(orgs, org)
	}

	return orgs, nil
}

func (s *SmartContract) CreateOrg(ctx contractapi.TransactionContextInterface, orgJson string) error {
	var org Org
	err := json.Unmarshal([]byte(orgJson), &org)
	if err != nil {
		return fmt.Errorf("failed to unmarshal org: %v", err)
	}

	did := org.Did

	result, _ := s.ReadOrg(ctx, did.ToString())
	if result.Did.SpecificID != "" {
		return fmt.Errorf("the org %s already exists", did)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(orgJson))
	if err != nil {
		return fmt.Errorf("failed to put org: %v", err)
	}

	return nil
}

func (s *SmartContract) UpdateOrg(ctx contractapi.TransactionContextInterface, orgJson string) error {
	var org Org
	err := json.Unmarshal([]byte(orgJson), &org)
	if err != nil {
		return fmt.Errorf("failed to unmarshal org: %v", err)
	}

	did := org.Did

	readOrg, err := s.ReadOrg(ctx, did.ToString())
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}
	if readOrg.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", did)
	}

	if org.Name != "" {
		readOrg.Name = org.Name
	}

	readJson, err := json.Marshal(readOrg)
	if err != nil {
		return fmt.Errorf("failed to marshal org: %v", err)
	}

	err = ctx.GetStub().PutState(did.ToString(), []byte(readJson))
	if err != nil {
		return fmt.Errorf("failed to put org: %v", err)
	}

	return nil
}

func (s *SmartContract) DeleteOrg(ctx contractapi.TransactionContextInterface, did string) error {
	readOrg, err := s.ReadOrg(ctx, did)
	if err != nil {
		return fmt.Errorf("failed to get org info: %v", err)
	}
	if readOrg.Did.SpecificID == "" {
		return fmt.Errorf("the org %s does not exist", did)
	}

	return ctx.GetStub().DelState(did)
}
